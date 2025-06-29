package internal

import (
	"encoding/json"
	"time"
)

type Webhook struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	Name      string
	Secret    string // HMAC key used to verify events
	Source    string // "github", "auth0", etc
	CreatedAt time.Time
}

type NotificationChannel struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	Name      string
	Type      string          // "slack", "webhook"
	Config    json.RawMessage // contains token/webhook/channel, etc
	CreatedAt time.Time
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
