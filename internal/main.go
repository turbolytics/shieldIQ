package internal

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Webhook struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	Name      string    `json:"name"`
	Secret    string    `json:"secret"` // HMAC key used to verify events
	Source    string    `json:"source"` // "github", "auth0", etc
	CreatedAt time.Time `json:"created_at"`
	Events    []string  `json:"events"`
}

type NotificationChannel struct {
	ID        uuid.UUID       `json:"id"`
	TenantID  uuid.UUID       `json:"tenant_id"`
	Name      string          `json:"name"`
	Type      string          `json:"type"`   // "slack", "webhook"
	Config    json.RawMessage `json:"config"` // contains token/webhook/channel, etc
	CreatedAt time.Time       `json:"created_at"`
}

type EvaluationType string

const (
	EvaluationTypeLiveTrigger EvaluationType = "LIVE_TRIGGER"
)

type AlertLevel string

const (
	AlertLevelLow    AlertLevel = "LOW"
	AlertLevelMedium AlertLevel = "MEDIUM"
	AlertLevelHigh   AlertLevel = "HIGH"
)

type Rule struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	EvaluationType EvaluationType `json:"evaluation_type"` // Only LIVE_TRIGGER for now
	EventSource    string         `json:"event_source"`    // e.g. "github", "auth0"
	EventType      string         `json:"event_type"`      // e.g. "github.pull_request"
	SQL            string         `json:"sql"`             // SQL expression, e.g. "reviewed = false"
	// NotificationChannels []string       `json:"notification_ids"` // Channels to alert
	CreatedAt  time.Time  `json:"created_at"`
	AlertLevel AlertLevel `json:"alert_level"`
	Active     bool       `json:"active"` // Whether the rule is currently active
}
