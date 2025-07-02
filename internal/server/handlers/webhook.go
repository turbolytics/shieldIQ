package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/turbolytics/sqlsec/internal/auth"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/turbolytics/sqlsec/internal"
	"github.com/turbolytics/sqlsec/internal/db"
)

type CreateWebhookRequest struct {
	Name   string   `json:"name"`
	Source string   `json:"source"`
	Events []string `json:"events"`
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

	// Hardcoded tenant_id for now
	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000000")

	id := uuid.New()
	createdAt := time.Now().UTC()
	secret, err := auth.GenerateSecret()
	if err != nil {
		wh.logger.Error("failed to generate secret", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	hook, err := wh.queries.CreateWebhook(r.Context(), db.CreateWebhookParams{
		ID:        id,
		TenantID:  tenantID,
		Name:      req.Name,
		Secret:    secret,
		Source:    req.Source,
		Events:    mustMarshalEvents(req.Events),
		CreatedAt: sql.NullTime{Time: createdAt, Valid: true},
	})
	if err != nil {
		wh.logger.Error("failed to create webhook", zap.Error(err))
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
		Events:    mustUnmarshalEvents(hook.Events),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (wh *Webhook) Get(w http.ResponseWriter, r *http.Request) {
	ctx := chi.RouteContext(r.Context())
	fmt.Println("URL Path:", r.URL.Path)
	fmt.Println(ctx)
	wh.logger.Info(
		"Get webhook called",
		zap.String("webhook_id", chi.URLParam(r, "webhook_id")),
		zap.Any("request", r),
	)
	idStr := chi.URLParam(r, "webhook_id")
	id, err := uuid.Parse(idStr)
	fmt.Println("here", idStr, "test")
	if err != nil {
		http.Error(w, "invalid webhook id", http.StatusBadRequest)
		return
	}
	// Hardcoded tenant_id for now
	tid := uuid.MustParse("00000000-0000-0000-0000-000000000000")

	webhook, err := wh.queries.GetWebhook(r.Context(), db.GetWebhookParams{
		ID:       id,
		TenantID: tid,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "webhook not found", http.StatusNotFound)
			return
		}
		wh.logger.Error("failed to get webhook", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(internal.Webhook{
		ID:        webhook.ID,
		TenantID:  webhook.TenantID,
		Name:      webhook.Name,
		Secret:    webhook.Secret,
		Source:    webhook.Source,
		CreatedAt: webhook.CreatedAt.Time,
	})
}

func mustMarshalEvents(events []string) []byte {
	data, err := json.Marshal(events)
	if err != nil {
		panic(err)
	}
	return data
}

func mustUnmarshalEvents(data []byte) []string {
	var events []string
	if err := json.Unmarshal(data, &events); err != nil {
		panic(err)
	}
	return events
}
