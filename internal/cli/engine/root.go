package engine

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/apache/arrow-adbc/go/adbc/drivermgr"
	"github.com/turbolytics/shieldIQ/internal/db/queries/alerts"
	"github.com/turbolytics/shieldIQ/internal/db/queries/events"
	"github.com/turbolytics/shieldIQ/internal/db/queries/rules"
	"go.uber.org/zap"
	"log"
	"os"

	"github.com/spf13/cobra"
	enginepkg "github.com/turbolytics/shieldIQ/internal/engine"
)

func NewCmd() *cobra.Command {
	var dsn string
	cmd := &cobra.Command{
		Use:   "engine",
		Short: "Engine related commands",
	}

	cmd.PersistentFlags().StringVar(&dsn, "dsn", os.Getenv("SHIELDIQ_DB_DSN"), "Postgres DSN (or set SHIELDIQ_DB_DSN env var)")
	cmd.AddCommand(newRunCmd(&dsn))
	cmd.AddCommand(newTestCmd())

	return cmd
}

func newRunCmd(dsn *string) *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Start the engine in daemon mode",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger, err := zap.NewDevelopment()
			if err != nil {
				log.Fatalf("failed to initialize zap logger: %v", err)
			}
			defer logger.Sync()

			var drv drivermgr.Driver
			db, err := drv.NewDatabase(map[string]string{
				"driver":     "duckdb",
				"entrypoint": "duckdb_adbc_init",
				"path":       ":memory:",
			})
			if err != nil {
				return fmt.Errorf("failed to initialize DuckDB driver: %w", err)
			}

			conn, err := db.Open(context.Background())
			if err != nil {
				return fmt.Errorf("failed to open DuckDB connection: %w", err)
			}
			defer func() {
				if err := conn.Close(); err != nil {
					logger.Error("failed to close DuckDB connection", zap.Error(err))
				}
			}()

			if *dsn == "" {
				log.Fatal("Postgres DSN must be set via --dsn or SHIELDIQ_DB_DSN env var")
			}
			dbConn, err := sql.Open("postgres", *dsn)
			if err != nil {
				log.Fatalf("failed to connect to database: %v", err)
			}
			defer dbConn.Close()
			eventQueries := events.New(dbConn)
			ruleQueries := rules.New(dbConn)
			alertQ := alerts.New(dbConn)
			alerter := enginepkg.NewPostgresAlerter(dbConn, alertQ, logger)
			engine := enginepkg.New(
				conn,
				eventQueries,
				ruleQueries,
				alerter,
				logger,
			)
			logger.Info("Starting engine daemon...")
			if err := engine.Run(context.Background()); err != nil {
				log.Fatalf("engine exited with error: %v", err)
			}
			return nil
		},
	}
}

func newTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "Test a rule execution against a payload",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("engine test called")
		},
	}
}
