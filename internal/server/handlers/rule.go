package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/turbolytics/sqlsec/internal"
	"github.com/turbolytics/sqlsec/internal/db"
	"github.com/turbolytics/sqlsec/internal/source"
	"net/http"
	"strconv"
	"time"
)

type RuleHandlers struct {
	queries *db.Queries
}

func NewRuleHandlers(queries *db.Queries) *RuleHandlers {
	return &RuleHandlers{queries: queries}
}

type RuleCreateRequest struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	Source         string `json:"source"`
	EventType      string `json:"event_type"`
	Condition      string `json:"condition"`
	EvaluationType string `json:"evaluation_type"`
	AlertLevel     string `json:"alert_level"`
}

type TestRuleRequest struct {
	Event map[string]interface{} `json:"event"`
}

type TestRuleResponse struct {
	Match      bool                   `json:"match"`
	AlertLevel string                 `json:"alert_level"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

func (h *RuleHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var req RuleCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	// Validate source
	if !source.DefaultRegistry.IsEnabled(req.Source) {
		http.Error(w, "unsupported source", http.StatusBadRequest)
		return
	}
	id := uuid.New()
	// TODO: get tenant_id from context/session
	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	createdAt := time.Now().UTC()
	// Insert into DB
	dbRule, err := h.queries.CreateRule(r.Context(), db.CreateRuleParams{
		ID:          id,
		TenantID:    tenantID,
		Name:        req.Name,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		Source:      req.Source,
		EventType:   req.EventType,
		Sql:         req.Condition,
		EvalType:    req.EvaluationType,
		AlertLevel:  req.AlertLevel,
		CreatedAt:   createdAt,
	})
	if err != nil {
		http.Error(w, "failed to create rule", http.StatusInternalServerError)
		return
	}
	resp := internal.Rule{
		ID:             dbRule.ID.String(),
		Name:           dbRule.Name,
		Description:    dbRule.Description.String,
		EvaluationType: internal.EvaluationType(dbRule.EvalType),
		EventSource:    dbRule.Source,
		EventType:      dbRule.EventType,
		SQL:            dbRule.Sql,
		CreatedAt:      dbRule.CreatedAt,
		AlertLevel:     internal.AlertLevel(dbRule.AlertLevel),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *RuleHandlers) List(w http.ResponseWriter, r *http.Request) {
	// TODO: get tenant_id from context/session
	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	limit := int32(50)
	offset := int32(0)
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 100 {
			limit = int32(v)
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = int32(v)
		}
	}
	rules, err := h.queries.ListRules(r.Context(), db.ListRulesParams{
		TenantID: tenantID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		http.Error(w, "failed to list rules", http.StatusInternalServerError)
		return
	}
	resp := make([]map[string]interface{}, 0, len(rules))
	for _, rule := range rules {
		resp = append(resp, map[string]interface{}{
			"id":              rule.ID,
			"tenant_id":       rule.TenantID,
			"name":            rule.Name,
			"description":     rule.Description.String,
			"source":          rule.Source,
			"event_type":      rule.EventType,
			"condition":       rule.Sql,
			"evaluation_type": rule.EvalType,
			"alert_level":     rule.AlertLevel,
			"created_at":      rule.CreatedAt,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *RuleHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	// TODO: get tenant_id from context/session
	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	err = h.queries.DeleteRule(r.Context(), db.DeleteRuleParams{
		ID:       id,
		TenantID: tenantID,
	})
	if err != nil {
		http.Error(w, "rule not found or failed to delete", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *RuleHandlers) Test(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement this
	idStr := chi.URLParam(r, "id")
	_, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var req TestRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	// TODO: fetch rule, evaluate condition, return result
	resp := TestRuleResponse{
		Match:      true,
		AlertLevel: "MEDIUM",
		Details:    map[string]interface{}{"info": "stubbed"},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *RuleHandlers) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	// TODO: get tenant_id from context/session
	tenantID := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	dbRule, err := h.queries.GetRuleByID(r.Context(), db.GetRuleByIDParams{
		ID:       id,
		TenantID: tenantID,
	})
	if err != nil {
		http.Error(w, "rule not found", http.StatusNotFound)
		return
	}
	resp := internal.Rule{
		ID:             dbRule.ID.String(),
		Name:           dbRule.Name,
		Description:    dbRule.Description.String,
		EvaluationType: internal.EvaluationType(dbRule.EvalType),
		EventSource:    dbRule.Source,
		EventType:      dbRule.EventType,
		SQL:            dbRule.Sql,
		CreatedAt:      dbRule.CreatedAt,
		AlertLevel:     internal.AlertLevel(dbRule.AlertLevel),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
