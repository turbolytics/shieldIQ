package alerter

import (
	"github.com/spf13/cobra"
)

func NewRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Run the alerter service",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement alerter run logic
			cmd.Println("Alerter run executed")
		},
	}
}
