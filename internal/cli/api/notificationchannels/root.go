package notificationchannels

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "channels",
		Short: "Manage notification channels",
	}

	cmd.AddCommand(NewCreateCmd())
	cmd.AddCommand(NewTestCmd())
	cmd.AddCommand(NewListCmd())

	return cmd
}
