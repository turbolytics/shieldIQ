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

type Rule struct {
	ID                    uuid.UUID
	TenantID              uuid.UUID
	Name                  string
	SQL                   string
	Severity              string // "low", "medium", "high"
	NotificationChannelID uuid.UUID
	Enabled               bool
	CreatedAt             time.Time
}

type Alert struct {
	ID        uuid.UUID
	RuleID    uuid.UUID
	TenantID  uuid.UUID
	Message   string
	Severity  string
	Data      json.RawMessage
	CreatedAt time.Time
}
