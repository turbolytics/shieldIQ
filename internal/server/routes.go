package server

import (
	"github.com/go-chi/chi"
	"github.com/turbolytics/sqlsec/internal/server/handlers"
)

func New(wh *handlers.Webhook) *chi.Mux {
	r := chi.NewRouter()

	// Middleware: logging, recover, etc.
	// r.Use(middleware.Logger)

	// Register routes
	r.Route("/api", func(r chi.Router) {
		r.Post("/webhooks", wh.Create)
	})

	return r
}
