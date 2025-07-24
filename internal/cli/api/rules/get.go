package rules

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

func NewGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <rule-id>",
		Short: "Get a rule by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ruleID := args[0]
			baseURL, _ := cmd.Flags().GetString("base-url")
			url := fmt.Sprintf("%s/api/rules/%s", baseURL, ruleID)
			resp, err := http.Get(url)
			if err != nil {
				return fmt.Errorf("failed to call get endpoint: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("server returned status: %s", resp.Status)
			}
			var rule map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&rule); err != nil {
				return fmt.Errorf("failed to decode response: %w", err)
			}
			t := table.NewWriter()
			t.SetOutputMirror(cmd.OutOrStdout())
			t.AppendHeader(table.Row{"Attribute", "Value"})
			for k, v := range rule {
				t.AppendRow(table.Row{k, v})
			}
			t.SetStyle(table.StyleDefault)
			t.Style().Options.SeparateRows = false
			t.Style().Box = table.StyleBoxDefault
			t.Style().Format.Header = text.FormatDefault
			t.Render()
			return nil
		},
	}
	return cmd
}
