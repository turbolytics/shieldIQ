package api

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/turbolytics/sqlsec/internal/cli/api/notificationchannels"
	"github.com/turbolytics/sqlsec/internal/cli/api/rules"
	"github.com/turbolytics/sqlsec/internal/cli/api/webhooks"
	"os"
)

var baseURL string

func NewCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "api",
		Short: "Manage api resources",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if baseURL == "" {
				return fmt.Errorf("missing required --base-url flag or SQLSEC_API_BASE_URL environment variable")
			}
			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&baseURL, "base-url", "", "Base server URL for API requests (or set SQLSEC_API_BASE_URL)")
	cobra.OnInitialize(initBaseURL)

	cmd.AddCommand(webhooks.NewCommand())
	cmd.AddCommand(notificationchannels.NewCommand())
	cmd.AddCommand(rules.NewCommand())

	return cmd
}

func initBaseURL() {
	if baseURL == "" {
		baseURL = os.Getenv("SQLSEC_API_BASE_URL")
	}
}
