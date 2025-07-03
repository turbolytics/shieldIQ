package handlers

import (
	"encoding/json"
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
		Name:      req.Name, // Name is not in DB yet, so use request value
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
	// TODO: parse id from URL, get channel, send test notification
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
