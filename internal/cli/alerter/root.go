package alerter

import (
	"github.com/spf13/cobra"
)

func NewAlerterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alerter",
		Short: "Alerter related commands",
	}
	cmd.AddCommand(NewRunCmd())
	cmd.AddCommand(NewTestCmd())
	return cmd
}
