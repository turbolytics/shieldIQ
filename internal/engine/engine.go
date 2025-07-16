package engine

import (
	"context"
	"database/sql"
	"errors"
	"github.com/apache/arrow-adbc/go/adbc"
	"github.com/turbolytics/sqlsec/internal/db/queries/events"
	"github.com/turbolytics/sqlsec/internal/db/queries/rules"
	"github.com/turbolytics/sqlsec/internal/engine/sandbox"
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

	box, err := sandbox.New(
		ctx,
		sandbox.WithDuckDBMemoryConnection(),
		sandbox.WithLogger(e.logger),
	)
	if err != nil {
		e.logger.Error("Failed to create sandbox",
			zap.Error(err),
		)
		return err
	}
	defer box.Close()

	if err := box.AddEvent(ctx, event); err != nil {
		e.logger.Error("Failed to add event to sandbox",
			zap.String("event_id", event.ID.String()),
			zap.Error(err),
		)
		return err
	}

	// 3. Loop through and run the rules against the event
	for _, rule := range rules {
		n, err := box.ExecuteRule(ctx, rule)
		if err != nil {
			// Log error, continue to next rule
			e.logger.Error("Failed to execute rule",
				zap.String("event_id", event.ID.String()),
				zap.String("rule_id", rule.ID.String()),
				zap.Error(err),
			)
			continue
		}
		// 4. Any rule execution that returns > 0 result should be saved to alerts table
		if n > 0 {
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
					zap.Int("result_count", n),
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

// saveAlert saves an alert to the alerts table if a rule is triggered.
func (e *Engine) saveAlert(ctx context.Context, rule rules.Rule, event *events.Event) error {
	// TODO: Implement saving alert to alerts table
	return nil
}

/*
The engine is the rules engine, responsible for processing events and evaluating rules.
*/
