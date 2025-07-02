package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/turbolytics/sqlsec/internal/server/handlers"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func RegisterRoutes(wr *chi.Mux, wh *handlers.Webhook, logger *zap.Logger) {
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

	wr.Use(logMiddleware)

	wr.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is healthy!"))
	})

	// Register routes
	wr.Route("/api", func(wr chi.Router) {
		wr.Route("/webhooks", func(wr chi.Router) {
			wr.Get("/{id}", wh.Get)
			wr.Post("/create", wh.Create)
		})
	})
}
