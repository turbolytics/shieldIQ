package engine

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "engine",
		Short: "Engine related commands",
	}

	cmd.AddCommand(newRunCmd())
	cmd.AddCommand(newTestCmd())

	return cmd
}

func newRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Start the engine in daemon mode",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("engine run called")
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
