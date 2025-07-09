package cli

import (
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/turbolytics/sqlsec/internal/cli/engine"
	"github.com/turbolytics/sqlsec/internal/cli/serve"
	"github.com/turbolytics/sqlsec/internal/cli/webhook"
	"log"
)

func NewRootCommand() *cobra.Command {
	/*
		connStr := os.Getenv("SQLSEC_DB_DSN")
		if connStr == "" {
			log.Fatal("SQLSEC_DB_DSN environment variable not set")
		}

		conn, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
		defer conn.Close()

		q := db.New(conn)
	*/

	rootCmd := &cobra.Command{
		Use:   "sqlsec",
		Short: "SQLSec – Security Alerting for SaaS Webhooks",
	}

	// Register subcommands
	rootCmd.AddCommand(webhook.NewCommand())
	rootCmd.AddCommand(serve.NewCommand())
	rootCmd.AddCommand(engine.NewCmd())

	return rootCmd
}

func Execute() {
	rootCmd := NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("error: %v", err)
	}
}
