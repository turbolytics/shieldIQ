package api

import (
	"github.com/spf13/cobra"
	"github.com/turbolytics/sqlsec/internal/cli/api/notificationchannels"
	"github.com/turbolytics/sqlsec/internal/cli/api/webhooks"
	"os"
)

func NewCommand() *cobra.Command {
	var baseURL string

	cmd := &cobra.Command{
		Use:   "api",
		Short: "Manage api resources",
	}

	cmd.PersistentFlags().StringVar(
		&baseURL,
		"base-url",
		os.Getenv("SQLSEC_API_BASE_URL"),
		"Base server URL for API requests",
	)

	cmd.AddCommand(webhooks.NewCommand())
	cmd.AddCommand(notificationchannels.NewCommand())

	return cmd
}
