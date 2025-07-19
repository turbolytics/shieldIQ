package engine

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/turbolytics/sqlsec/internal/db/queries/alerts"
	"github.com/turbolytics/sqlsec/internal/db/queries/events"
	"github.com/turbolytics/sqlsec/internal/db/queries/rules"
	"go.uber.org/zap"
	"testing"
)

// Stub implementations for Querier interfaces

type mockAlertsQuerier struct {
	called bool
}

func (m *mockAlertsQuerier) CreateAlert(ctx context.Context, arg alerts.CreateAlertParams) (alerts.CreateAlertRow, error) {
	m.called = true
	return alerts.CreateAlertRow{}, nil
}
func (m *mockAlertsQuerier) FetchNextAlertForProcessing(ctx context.Context, lockedBy sql.NullString) (uuid.UUID, error) {
	return uuid.New(), nil
}
func (m *mockAlertsQuerier) InsertAlertProcessingQueue(ctx context.Context, alertID uuid.UUID) (alerts.AlertProcessingQueue, error) {
	return alerts.AlertProcessingQueue{}, nil
}
func (m *mockAlertsQuerier) MarkAlertProcessingDelivered(ctx context.Context, alertID uuid.UUID) error {
	return nil
}
func (m *mockAlertsQuerier) MarkAlertProcessingFailed(ctx context.Context, arg alerts.MarkAlertProcessingFailedParams) error {
	return nil
}

// alerts.Querier only has CreateAlert

// --- events.Querier ---
type mockEventsQuerier struct{}

func (m *mockEventsQuerier) FetchNextEventForProcessing(ctx context.Context, lockedBy sql.NullString) (uuid.UUID, error) {
	return uuid.New(), nil
}
func (m *mockEventsQuerier) GetEventByID(ctx context.Context, id uuid.UUID) (events.Event, error) {
	return events.Event{
		ID:        id,
		EventType: "test_event",
		Source:    "test_source",
		TenantID:  uuid.New(),
	}, nil
}
func (m *mockEventsQuerier) InsertEvent(ctx context.Context, arg events.InsertEventParams) (events.Event, error) {
	panic("not implemented")
}
func (m *mockEventsQuerier) InsertEventProcessingQueue(ctx context.Context, eventID uuid.UUID) (events.EventProcessingQueue, error) {
	panic("not implemented")
}
func (m *mockEventsQuerier) MarkEventProcessingDone(ctx context.Context, eventID uuid.UUID) error {
	return nil
}
func (m *mockEventsQuerier) MarkEventProcessingFailed(ctx context.Context, arg events.MarkEventProcessingFailedParams) error {
	panic("not implemented")
}

// --- rules.Querier ---
type mockRulesQuerier struct{}

func (m *mockRulesQuerier) CreateRule(ctx context.Context, arg rules.CreateRuleParams) (rules.Rule, error) {
	panic("not implemented")
}
func (m *mockRulesQuerier) CreateRuleDestination(ctx context.Context, arg rules.CreateRuleDestinationParams) (rules.RuleDestination, error) {
	panic("not implemented")
}
func (m *mockRulesQuerier) DeleteRule(ctx context.Context, arg rules.DeleteRuleParams) error {
	panic("not implemented")
}
func (m *mockRulesQuerier) DeleteRuleDestination(ctx context.Context, arg rules.DeleteRuleDestinationParams) error {
	panic("not implemented")
}
func (m *mockRulesQuerier) GetRuleByID(ctx context.Context, arg rules.GetRuleByIDParams) (rules.Rule, error) {
	panic("not implemented")
}
func (m *mockRulesQuerier) GetRulesForEvent(ctx context.Context, arg rules.GetRulesForEventParams) ([]rules.Rule, error) {
	return []rules.Rule{
		{
			ID:       uuid.New(),
			Name:     "test_rule",
			TenantID: arg.TenantID,           // associate with tenant
			Active:   true,                   // rule is enabled
			Sql:      "SELECT * FROM events", // returns everything
		},
	}, nil
}
func (m *mockRulesQuerier) ListNotificationChannelsForRule(ctx context.Context, ruleID uuid.UUID) ([]rules.NotificationChannel, error) {
	panic("not implemented")
}
func (m *mockRulesQuerier) ListRuleDestinationChannelIDs(ctx context.Context, ruleID uuid.UUID) ([]uuid.UUID, error) {
	panic("not implemented")
}
func (m *mockRulesQuerier) ListRules(ctx context.Context, arg rules.ListRulesParams) ([]rules.Rule, error) {
	panic("not implemented")
}
func (m *mockRulesQuerier) UpdateRuleActive(ctx context.Context, arg rules.UpdateRuleActiveParams) (rules.Rule, error) {
	panic("not implemented")
}

// --- mockAlerter for Alerter interface ---
type mockAlerter struct {
	called        bool
	receivedRule  rules.Rule
	receivedEvent *events.Event
}

func (m *mockAlerter) CreateAlert(ctx context.Context, rule rules.Rule, event *events.Event) error {
	m.called = true
	m.receivedRule = rule
	m.receivedEvent = event
	return nil
}

func TestEngine_ExecuteOnce_HappyPath(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop()
	alerter := &mockAlerter{}
	engine := &Engine{
		alerter:      alerter,
		eventQueries: &mockEventsQuerier{},
		ruleQueries:  &mockRulesQuerier{},
		logger:       logger,
	}

	err := engine.ExecuteOnce(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !alerter.called {
		t.Error("expected CreateAlert to be called, but it was not")
	}
	if alerter.receivedRule.Name != "test_rule" {
		t.Errorf("expected rule name 'test_rule', got '%s'", alerter.receivedRule.Name)
	}
	if alerter.receivedEvent == nil || alerter.receivedEvent.EventType != "test_event" {
		t.Errorf("expected event type 'test_event', got '%v'", alerter.receivedEvent)
	}
}
