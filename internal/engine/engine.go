package engine

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/apache/arrow-adbc/go/adbc"
	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/turbolytics/sqlsec/internal/db/queries/events"
	"github.com/turbolytics/sqlsec/internal/db/queries/rules"
	"go.uber.org/zap"
	"time"
)

// Engine is responsible for processing events and evaluating rules.
type Engine struct {
	conn adbc.Connection

	eventQueries events.Querier
	ruleQueries  rules.Querier
	logger       *zap.Logger
}

func New(conn adbc.Connection, eventQueries *events.Queries, ruleQueries *rules.Queries, l *zap.Logger) *Engine {
	return &Engine{
		conn:         conn,
		eventQueries: eventQueries,
		ruleQueries:  ruleQueries,
		logger:       l,
	}
}

// ExecuteOnce is the main execution entrypoint for a single engine loop.
func (e *Engine) ExecuteOnce(ctx context.Context) error {
	// 1. Fetch an event to process (from event_processing_queue)
	event, err := e.fetchNextEvent(ctx)
	if err != nil {
		return err
	}
	if event == nil {
		// No event to process
		return nil
	}

	e.logger.Debug("Processing event",
		zap.String("event_id", event.ID.String()),
		zap.String("event_type", event.EventType),
		zap.String("source", event.Source),
		zap.String("tenant_id", event.TenantID.String()),
		zap.Duration("lag", time.Since(event.ReceivedAt.Time)),
	)

	// 2. Fetch rules associated with that event type
	rules, err := e.fetchRulesForEvent(ctx, event)
	if err != nil {
		return err
	}
	e.logger.Debug("Fetched rules for event",
		zap.String("event_id", event.ID.String()),
		zap.Int("rule_count", len(rules)),
	)

	// 3. Loop through and run the rules against the event
	for _, rule := range rules {
		result, err := e.executeRule(ctx, rule, event)
		if err != nil {
			// Log error, continue to next rule
			continue
		}
		// 4. Any rule execution that returns > 0 result should be saved to alerts table
		if result > 0 {
			if err := e.saveAlert(ctx, rule, event); err != nil {
				e.logger.Error("Failed to save alert",
					zap.String("event_id", event.ID.String()),
					zap.String("rule_id", rule.ID.String()),
					zap.Error(err),
				)
			} else {
				e.logger.Debug("Alert saved",
					zap.String("event_id", event.ID.String()),
					zap.String("rule_id", rule.ID.String()),
					zap.Int("result_count", result),
				)
			}
		}
	}
	// 4. Mark the event as processed
	err = e.eventQueries.MarkEventProcessingDone(ctx, event.ID)

	return err
}

// Run starts the engine in daemon mode.
func (e *Engine) Run(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			e.ExecuteOnce(ctx)
		}
	}
}

// fetchNextEvent fetches the next event to process from the queue.
func (e *Engine) fetchNextEvent(ctx context.Context) (*events.Event, error) {
	// Use the FetchNextEventForProcessing query to get the next event_id
	lockedBy := sql.NullString{String: "engine-daemon", Valid: true}
	eventID, err := e.eventQueries.FetchNextEventForProcessing(ctx, lockedBy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No event to process
		}
		return nil, err
	}
	event, err := e.eventQueries.GetEventByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// fetchRulesForEvent fetches rules associated with the event type.
func (e *Engine) fetchRulesForEvent(ctx context.Context, event *events.Event) ([]rules.Rule, error) {
	return e.ruleQueries.GetRulesForEvent(ctx, rules.GetRulesForEventParams{
		TenantID:  event.TenantID,
		Source:    event.Source,
		EventType: event.EventType,
	})
}

// executeRule runs a rule against an event and returns the result count.
func (e *Engine) executeRule(ctx context.Context, rule rules.Rule, event *events.Event) (int, error) {
	if rule.Sql == "" {
		return 0, fmt.Errorf("rule SQL is empty for rule ID %s", rule.ID)
	}

	// 1. Create the table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS events (
		id UUID,
		tenant_id UUID,
		webhook_id UUID,
		source TEXT,
		event_type TEXT,
		action TEXT,
		raw_payload JSON,
		dedup_hash TEXT,
		received_at TIMESTAMP
	);
	TRUNCATE TABLE events;`

	createTableStmt, err := e.conn.NewStatement()
	if err != nil {
		return 0, fmt.Errorf("create statement: %w", err)
	}
	defer createTableStmt.Close()
	if err := createTableStmt.SetSqlQuery(createTableSQL); err != nil {
		return -1, fmt.Errorf("set query error: %v", err)
	}
	reader, _, err := createTableStmt.ExecuteQuery(ctx)
	if err != nil {
		return -1, fmt.Errorf("query execution error: %v", err)
	}
	defer reader.Release()
	for reader.Next() {
		reader.Release()
	}

	// 2. Insert the event
	insertSQL := `INSERT INTO events (
		id, tenant_id, webhook_id, source, event_type, action, raw_payload, dedup_hash, received_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	insertStmt, err := e.conn.NewStatement()
	if err != nil {
		return -1, fmt.Errorf("new insert stmt: %w", err)
	}
	defer insertStmt.Close()

	if err := insertStmt.SetSqlQuery(insertSQL); err != nil {
		return -1, fmt.Errorf("set insert sql: %w", err)
	}

	r := arrowRecordFromEvent(event)
	if err := insertStmt.Bind(ctx, r); err != nil {
		return -1, fmt.Errorf("bind insert: %w", err)
	}

	if _, err := insertStmt.ExecuteUpdate(ctx); err != nil {
		return -1, fmt.Errorf("execute insert: %w", err)
	}

	// 3. Execute the rule SQL
	ruleStmt, err := e.conn.NewStatement()
	if err != nil {
		return -1, fmt.Errorf("new rule stmt: %w", err)
	}
	defer ruleStmt.Close()

	e.logger.Debug("Executing rule SQL",
		zap.String("rule_id", rule.ID.String()),
	)

	if err := ruleStmt.SetSqlQuery(rule.Sql); err != nil {
		return -1, fmt.Errorf("set rule sql: %w", err)
	}
	reader, _, err = ruleStmt.ExecuteQuery(ctx)
	if err != nil {
		return -1, fmt.Errorf("execute rule sql: %w", err)
	}
	defer reader.Release()

	// 4. Count rows
	count := 0
	for reader.Next() {
		count++
	}
	if err := reader.Err(); err != nil {
		return -1, fmt.Errorf("read rule results: %w", err)
	}

	return count, nil
}

// saveAlert saves an alert to the alerts table if a rule is triggered.
func (e *Engine) saveAlert(ctx context.Context, rule rules.Rule, event *events.Event) error {
	// TODO: Implement saving alert to alerts table
	return nil
}

/*
The engine is the rules engine, responsible for processing events and evaluating rules.
*/

func arrowRecordFromEvent(event *events.Event) arrow.Record {
	pool := memory.NewGoAllocator()

	schema := arrow.NewSchema([]arrow.Field{
		{Name: "id", Type: arrow.BinaryTypes.String},
		{Name: "tenant_id", Type: arrow.BinaryTypes.String},
		{Name: "webhook_id", Type: arrow.BinaryTypes.String},
		{Name: "source", Type: arrow.BinaryTypes.String},
		{Name: "event_type", Type: arrow.BinaryTypes.String},
		{Name: "action", Type: arrow.BinaryTypes.String},
		{Name: "raw_payload", Type: arrow.BinaryTypes.String},
		{Name: "dedup_hash", Type: arrow.BinaryTypes.String},
		{Name: "received_at", Type: arrow.FixedWidthTypes.Timestamp_ms},
	}, nil)

	b := array.NewRecordBuilder(pool, schema)
	defer b.Release()

	payload, _ := event.RawPayload.MarshalJSON()

	b.Field(0).(*array.StringBuilder).Append(event.ID.String())
	b.Field(1).(*array.StringBuilder).Append(event.TenantID.String())
	b.Field(2).(*array.StringBuilder).Append(event.WebhookID.String())
	b.Field(3).(*array.StringBuilder).Append(event.Source)
	b.Field(4).(*array.StringBuilder).Append(event.EventType)
	b.Field(5).(*array.StringBuilder).Append(event.Action.String)
	b.Field(6).(*array.StringBuilder).Append(string(payload))
	b.Field(7).(*array.StringBuilder).Append(event.DedupHash.String)
	b.Field(8).(*array.TimestampBuilder).Append(arrow.Timestamp(event.ReceivedAt.Time.UnixMilli()))

	return b.NewRecord()
}
