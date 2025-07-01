package server

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/turbolytics/sqlsec/internal/server/handlers"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func New(wh *handlers.Webhook, logger *zap.Logger) *chi.Mux {
	r := chi.NewRouter()

	logMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				logger.Info("request",
					zap.String("from", r.RemoteAddr),
					zap.String("protocol", r.Proto),
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.Int("status", ww.Status()),
					zap.Int("bytes", ww.BytesWritten()),
					zap.Duration("duration", time.Since(start)),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}

	r.Use(logMiddleware)

	// Register routes
	r.Route("/api", func(r chi.Router) {
		r.Post("/webhooks", wh.Create)
		r.Get("/webhooks/{id}", wh.Get)
	})

	return r
}
