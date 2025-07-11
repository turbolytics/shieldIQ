package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/turbolytics/sqlsec/internal"
	"github.com/turbolytics/sqlsec/internal/db"
	"github.com/turbolytics/sqlsec/internal/notify"
	_ "github.com/turbolytics/sqlsec/internal/notify/slack"
	"net/http"
)

type NotificationHandlers struct {
	queries *db.Queries
}

func NewNotificationHandlers(queries *db.Queries) *NotificationHandlers {
	return &NotificationHandlers{queries: queries}
}

type CreateNotificationChannelRequest struct {
	Name   string                 `json:"name"`
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

func (h *NotificationHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateNotificationChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	id := uuid.New()

	// TODO: get tenant_id from context/session
	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000000")

	// Validate notification channel type
	if !notify.DefaultRegistry.IsEnabled(notify.ChannelType(req.Type)) {
		http.Error(w, "unsupported notification channel type", http.StatusBadRequest)
		return
	}

	configJSON, _ := json.Marshal(req.Config)
	ch, err := h.queries.CreateNotificationChannel(r.Context(), db.CreateNotificationChannelParams{
		ID:       id,
		TenantID: tenantID,
		Name:     req.Name,
		Type:     req.Type,
		Config:   configJSON,
	})
	if err != nil {
		http.Error(w, "failed to create", http.StatusInternalServerError)
		return
	}
	resp := internal.NotificationChannel{
		ID:        ch.ID,
		TenantID:  ch.TenantID,
		Name:      ch.Name,
		Type:      ch.Type,
		Config:    ch.Config,
		CreatedAt: ch.CreatedAt.Time, // ch.CreatedAt is sql.NullTime
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *NotificationHandlers) List(w http.ResponseWriter, r *http.Request) {
	// TODO: get tenant_id from context/session
	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000000")

	chs, err := h.queries.ListNotificationChannels(r.Context(), tenantID)
	if err != nil {
		http.Error(w, "failed to list", http.StatusInternalServerError)
		return
	}
	resp := make([]internal.NotificationChannel, 0, len(chs))
	for _, ch := range chs {
		resp = append(resp, internal.NotificationChannel{
			ID:        ch.ID,
			TenantID:  ch.TenantID,
			Name:      "", // Name not in DB, left blank
			Type:      ch.Type,
			Config:    ch.Config,
			CreatedAt: ch.CreatedAt.Time,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *NotificationHandlers) Test(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	ch, err := h.queries.GetNotificationChannelByID(r.Context(), id)
	if err != nil {
		http.Error(w, "notification channel not found", http.StatusNotFound)
		return
	}
	// Parse config JSON to map[string]string
	var cfg map[string]string
	if err := json.Unmarshal(ch.Config, &cfg); err != nil {
		http.Error(w, "invalid config", http.StatusBadRequest)
		return
	}
	n, err := notify.DefaultRegistry.Get(notify.ChannelType(ch.Type))
	if err != nil {
		http.Error(w, "unsupported notification channel type", http.StatusBadRequest)
		return
	}
	err = n.Test(r.Context(), cfg)
	if err != nil {
		http.Error(w, "failed to send test notification: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Test message sent"))
}
