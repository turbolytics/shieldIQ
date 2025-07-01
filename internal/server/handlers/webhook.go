package handlers

import (
	"database/sql"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/turbolytics/sqlsec/internal"
	"github.com/turbolytics/sqlsec/internal/db"
)

type CreateWebhookRequest struct {
	TenantID string `json:"tenant_id"`
	Name     string `json:"name"`
	Source   string `json:"source"`
	Secret   string `json:"secret"`
}

type Server struct {
	queries *db.Queries
}

// NewWebhook creates a new Webhook handler with the given options.
func NewWebhook(queries *db.Queries, logger *zap.Logger) *Webhook {
	return &Webhook{
		queries: queries,
		logger:  logger,
	}
}

type Webhook struct {
	queries *db.Queries
	logger  *zap.Logger
}

func (wh *Webhook) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	tenantID, err := uuid.Parse(req.TenantID)
	if err != nil {
		http.Error(w, "invalid tenant_id", http.StatusBadRequest)
		return
	}

	id := uuid.New()
	createdAt := time.Now().UTC()

	hook, err := wh.queries.CreateWebhook(r.Context(), db.CreateWebhookParams{
		ID:        id,
		TenantID:  tenantID,
		Name:      req.Name,
		Secret:    req.Secret,
		Source:    req.Source,
		CreatedAt: sql.NullTime{Time: createdAt, Valid: true},
	})
	if err != nil {
		http.Error(w, "failed to create webhook", http.StatusInternalServerError)
		return
	}

	resp := internal.Webhook{
		ID:        hook.ID,
		TenantID:  hook.TenantID,
		Name:      hook.Name,
		Secret:    hook.Secret,
		Source:    hook.Source,
		CreatedAt: hook.CreatedAt.Time,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
