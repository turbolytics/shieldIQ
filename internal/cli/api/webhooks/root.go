package webhooks

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "webhooks",
		Short: "Manage webhooks",
	}

	cmd.AddCommand(NewWebhookCreateCmd())

	return cmd
}
