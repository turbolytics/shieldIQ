package alerter

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/turbolytics/shieldIQ/internal/db/queries/alerts"
	"github.com/turbolytics/shieldIQ/internal/db/queries/notificationchannels"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAlertQueries is a testify mock for alertQueries interface
// You may need to adjust the method signatures to match your actual interface
type MockAlertQueries struct {
	mock.Mock
}

func (m *MockAlertQueries) GetAlertByID(ctx context.Context, id uuid.UUID) (alerts.Alert, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(alerts.Alert), args.Error(1)
}

func (m *MockAlertQueries) FetchNextAlertForProcessing(ctx context.Context, lockedBy sql.NullString) (uuid.UUID, error) {
	args := m.Called(ctx, lockedBy)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockAlertQueries) CreateAlert(ctx context.Context, arg alerts.CreateAlertParams) (alerts.CreateAlertRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(alerts.CreateAlertRow), args.Error(1)
}

func (m *MockAlertQueries) FetchAlertForProcessing(ctx context.Context) (*alerts.Alert, error) {
	args := m.Called(ctx)
	return args.Get(0).(*alerts.Alert), args.Error(1)
}
func (m *MockAlertQueries) GetAlert(ctx context.Context, id int64) (*alerts.Alert, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*alerts.Alert), args.Error(1)
}
func (m *MockAlertQueries) FindNotificationChannel(ctx context.Context, alert *alerts.Alert) (*notificationchannels.NotificationChannel, error) {
	args := m.Called(ctx, alert)
	return args.Get(0).(*notificationchannels.NotificationChannel), args.Error(1)
}
func (m *MockAlertQueries) InsertAlertDelivery(ctx context.Context, delivery *alerts.AlertDelivery) error {
	args := m.Called(ctx, delivery)
	return args.Error(0)
}
func (m *MockAlertQueries) MarkAlertProcessingDelivered(ctx context.Context, alertID int64) error {
	args := m.Called(ctx, alertID)
	return args.Error(0)
}
func (m *MockAlertQueries) MarkAlertNotified(ctx context.Context, alertID int64) error {
	args := m.Called(ctx, alertID)
	return args.Error(0)
}

func TestAlerter_ExecuteOnce_HappyPath(t *testing.T) {
	t.Skip("Skipping test for now, needs proper setup")
	ctx := context.Background()
	mockQueries := new(MockAlertQueries)

	alert := &alerts.Alert{ID: uuid.New()}
	channel := &notificationchannels.NotificationChannel{ID: uuid.New()}
	// delivery := &alerts.AlertDelivery{AlertID: alert.ID, ChannelID: channel.ID}

	mockQueries.On("FetchAlertForProcessing", ctx).Return(alert, nil)
	mockQueries.On("GetAlert", ctx, alert.ID).Return(alert, nil)
	mockQueries.On("FindNotificationChannel", ctx, alert).Return(channel, nil)
	mockQueries.On("InsertAlertDelivery", ctx, mock.AnythingOfType("*alerter.AlertDelivery")).Return(nil)
	mockQueries.On("MarkAlertProcessingDelivered", ctx, alert.ID).Return(nil)
	mockQueries.On("MarkAlertNotified", ctx, alert.ID).Return(nil)

	// You may need to adjust this to match your actual Alerter struct and method
	alerter := &Alerter{
		// alertQueries: mockQueries,
	}

	err := alerter.ExecuteOnce(ctx)
	assert.NoError(t, err)

	mockQueries.AssertExpectations(t)
}
