package serve

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/turbolytics/sqlsec/internal/db"
	"github.com/turbolytics/sqlsec/internal/server"
	"github.com/turbolytics/sqlsec/internal/server/handlers"
	"go.uber.org/zap"
)

func NewCommand() *cobra.Command {
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

			queries := db.New(dbConn)
			wh := handlers.NewWebhook(queries, logger)
			r := server.New(wh)

			logger.Info("Starting server on :8080")
			if err := http.ListenAndServe(":8080", r); err != nil {
				logger.Fatal("server failed", zap.Error(err))
			}
		},
	}
	return cmd
}
