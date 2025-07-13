package engine

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/turbolytics/sqlsec/internal/db/queries/events"
	"github.com/turbolytics/sqlsec/internal/db/queries/rules"
	"go.uber.org/zap"
	"log"
	"os"

	"github.com/spf13/cobra"
	enginepkg "github.com/turbolytics/sqlsec/internal/engine"
)

func NewCmd() *cobra.Command {
	var dsn string
	cmd := &cobra.Command{
		Use:   "engine",
		Short: "Engine related commands",
	}

	cmd.PersistentFlags().StringVar(&dsn, "dsn", os.Getenv("SQLSEC_DB_DSN"), "Postgres DSN (or set SQLSEC_DB_DSN env var)")
	cmd.AddCommand(newRunCmd(&dsn))
	cmd.AddCommand(newTestCmd())

	return cmd
}

func newRunCmd(dsn *string) *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Start the engine in daemon mode",
		Run: func(cmd *cobra.Command, args []string) {
			logger, err := zap.NewDevelopment()
			if err != nil {
				log.Fatalf("failed to initialize zap logger: %v", err)
			}
			defer logger.Sync()

			if *dsn == "" {
				log.Fatal("Postgres DSN must be set via --dsn or SQLSEC_DB_DSN env var")
			}
			dbConn, err := sql.Open("postgres", *dsn)
			if err != nil {
				log.Fatalf("failed to connect to database: %v", err)
			}
			defer dbConn.Close()
			eventQueries := events.New(dbConn)
			ruleQueries := rules.New(dbConn)
			engine := enginepkg.New(
				eventQueries,
				ruleQueries,
				logger,
			)
			logger.Info("Starting engine daemon...")
			if err := engine.Run(context.Background()); err != nil {
				log.Fatalf("engine exited with error: %v", err)
			}
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
