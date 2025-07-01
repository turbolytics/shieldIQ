package webhook

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewWebhookCreateCmd() *cobra.Command {
	var name string
	var source string
	var tenantID string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new webhook source",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("here!")
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Name of the webhook (required)")
	cmd.Flags().StringVar(&source, "source", "", "Source system (e.g. github, auth0)")
	cmd.Flags().StringVar(&tenantID, "tenant-id", "", "Tenant UUID (required)")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("source")
	cmd.MarkFlagRequired("tenant-id")

	return cmd
}
