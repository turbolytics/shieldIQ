package engine

import (
	"context"
	"database/sql"
	"github.com/apache/arrow-adbc/go/adbc/drivermgr"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"

	events "github.com/turbolytics/sqlsec/internal/db/queries/events"
	rules "github.com/turbolytics/sqlsec/internal/db/queries/rules"
)

type stubEventsQueries struct{}

func (s *stubEventsQueries) FetchNextEventForProcessing(ctx context.Context, lockedBy sql.NullString) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}
func (s *stubEventsQueries) GetEventByID(ctx context.Context, id uuid.UUID) (events.Event, error) {
	return events.Event{}, nil
}
func (s *stubEventsQueries) InsertEvent(ctx context.Context, arg events.InsertEventParams) (events.Event, error) {
	return events.Event{}, nil
}
func (s *stubEventsQueries) InsertEventProcessingQueue(ctx context.Context, eventID uuid.UUID) (events.EventProcessingQueue, error) {
	return events.EventProcessingQueue{}, nil
}
func (s *stubEventsQueries) MarkEventProcessingDone(ctx context.Context, eventID uuid.UUID) error {
	return nil
}
func (s *stubEventsQueries) MarkEventProcessingFailed(ctx context.Context, arg events.MarkEventProcessingFailedParams) error {
	return nil
}

type stubRulesQueries struct{}

func (s *stubRulesQueries) CreateRule(ctx context.Context, arg rules.CreateRuleParams) (rules.Rule, error) {
	return rules.Rule{}, nil
}
func (s *stubRulesQueries) CreateRuleDestination(ctx context.Context, arg rules.CreateRuleDestinationParams) (rules.RuleDestination, error) {
	return rules.RuleDestination{}, nil
}
func (s *stubRulesQueries) DeleteRule(ctx context.Context, arg rules.DeleteRuleParams) error {
	return nil
}
func (s *stubRulesQueries) DeleteRuleDestination(ctx context.Context, arg rules.DeleteRuleDestinationParams) error {
	return nil
}
func (s *stubRulesQueries) GetRuleByID(ctx context.Context, arg rules.GetRuleByIDParams) (rules.Rule, error) {
	return rules.Rule{}, nil
}
func (s *stubRulesQueries) GetRulesForEvent(ctx context.Context, arg rules.GetRulesForEventParams) ([]rules.Rule, error) {
	return []rules.Rule{}, nil
}
func (s *stubRulesQueries) ListNotificationChannelsForRule(ctx context.Context, ruleID uuid.UUID) ([]rules.NotificationChannel, error) {
	return []rules.NotificationChannel{}, nil
}
func (s *stubRulesQueries) ListRuleDestinationChannelIDs(ctx context.Context, ruleID uuid.UUID) ([]uuid.UUID, error) {
	return []uuid.UUID{}, nil
}
func (s *stubRulesQueries) ListRules(ctx context.Context, arg rules.ListRulesParams) ([]rules.Rule, error) {
	return []rules.Rule{}, nil
}

func (s *stubRulesQueries) ExecuteRule(ctx context.Context, rule rules.Rule, event *events.Event) (int, error) {
	return 1, nil
}

func TestEngineExecuteRule_SingleResult(t *testing.T) {
	ctx := context.Background()
	stubEvents := &stubEventsQueries{}
	stubRules := &stubRulesQueries{}

	// Create dummy rule and event
	rule := rules.Rule{
		Sql: "select * from events where raw_payload->>'key' = 'value'",
	}
	event := events.Event{
		ID:         uuid.New(),
		TenantID:   uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		Source:     "test_source",
		EventType:  "test_event",
		RawPayload: []byte(`{"key": "value"}`),
	}

	var drv drivermgr.Driver
	db, err := drv.NewDatabase(map[string]string{
		"driver":     "duckdb",
		"entrypoint": "duckdb_adbc_init",
		"path":       ":memory:",
	})
	assert.NoError(t, err)
	conn, err := db.Open(context.Background())
	assert.NoError(t, err)

	defer func() {
		if err := conn.Close(); err != nil {
			t.Fatalf("failed to close DuckDB connection: %v", err)
		}
	}()

	engine := &Engine{
		conn:         conn,
		eventQueries: stubEvents,
		ruleQueries:  stubRules,
		logger:       zap.NewNop(),
	}

	result, err := engine.executeRule(ctx, rule, &event)
	assert.NoError(t, err)
	assert.Equal(t, 1, result)
}

func TestEngineExecuteRule_NoMatch(t *testing.T) {
	ctx := context.Background()
	stubEvents := &stubEventsQueries{}
	stubRules := &stubRulesQueries{}

	// Create dummy rule and event
	rule := rules.Rule{
		Sql: "select * from events where raw_payload->>'key' = 'nomatch'",
	}
	event := events.Event{
		ID:         uuid.New(),
		TenantID:   uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		Source:     "test_source",
		EventType:  "test_event",
		RawPayload: []byte(`{"key": "value"}`),
	}

	var drv drivermgr.Driver
	db, err := drv.NewDatabase(map[string]string{
		"driver":     "duckdb",
		"entrypoint": "duckdb_adbc_init",
		"path":       ":memory:",
	})
	assert.NoError(t, err)
	conn, err := db.Open(context.Background())
	assert.NoError(t, err)

	defer func() {
		if err := conn.Close(); err != nil {
			t.Fatalf("failed to close DuckDB connection: %v", err)
		}
	}()

	engine := &Engine{
		conn:         conn,
		eventQueries: stubEvents,
		ruleQueries:  stubRules,
		logger:       zap.NewNop(),
	}

	result, err := engine.executeRule(ctx, rule, &event)
	assert.NoError(t, err)
	assert.Equal(t, 0, result)
}
