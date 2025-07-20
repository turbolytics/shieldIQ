package alerter

import (
	"github.com/spf13/cobra"
	"os"
)

func NewAlerterCmd() *cobra.Command {
	var dsn string
	cmd := &cobra.Command{
		Use:   "alerter",
		Short: "Alerter related commands",
	}
	cmd.PersistentFlags().StringVar(&dsn, "dsn", os.Getenv("SQLSEC_DB_DSN"), "Postgres DSN for database connection")
	cmd.AddCommand(NewRunCmd(&dsn))
	cmd.AddCommand(NewTestCmd())
	return cmd
}
