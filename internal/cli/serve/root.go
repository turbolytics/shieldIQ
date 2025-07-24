package serve

import (
	"database/sql"
	"github.com/turbolytics/sqlsec/internal/db/queries/events"
	"github.com/turbolytics/sqlsec/internal/db/queries/notificationchannels"
	"github.com/turbolytics/sqlsec/internal/db/queries/rules"
	"github.com/turbolytics/sqlsec/internal/db/queries/webhooks"
	"github.com/turbolytics/sqlsec/internal/server"
	"github.com/turbolytics/sqlsec/internal/server/handlers"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewCommand() *cobra.Command {
	var port string
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the SQLSec API server",
		Run: func(cmd *cobra.Command, args []string) {
			logger, err := zap.NewProduction()
			if err != nil {
				log.Fatalf("failed to initialize zap logger: %v", err)
			}
			defer logger.Sync()

			dsn := os.Getenv("SQLSEC_DB_DSN")
			if dsn == "" {
				logger.Fatal("SQLSEC_DB_DSN environment variable not set")
			}

			dbConn, err := sql.Open("postgres", dsn)
			if err != nil {
				logger.Fatal("failed to connect to database", zap.Error(err))
			}
			defer dbConn.Close()

			eventQueries := events.New(dbConn)
			ruleQueries := rules.New(dbConn)
			ncQueries := notificationchannels.New(dbConn)
			webhookQueries := webhooks.New(dbConn)
			wh := handlers.NewWebhook(
				dbConn,
				eventQueries,
				webhookQueries,
				logger,
			)
			nh := handlers.NewNotificationHandlers(logger, ncQueries)
			rh := handlers.NewRuleHandlers(logger, ruleQueries)
			dh := handlers.NewDestinationHandlers(ruleQueries, ncQueries)
			router := chi.NewRouter()
			server.RegisterRoutes(router, wh, nh, rh, dh, logger)

			addr := ":" + port
			logger.Info("Starting server", zap.String("addr", addr))
			if err := http.ListenAndServe(addr, router); err != nil {
				logger.Fatal("server failed", zap.Error(err))
			}
		},
	}
	cmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to listen on")
	return cmd
}
