package rules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new rule",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			description, _ := cmd.Flags().GetString("description")
			source, _ := cmd.Flags().GetString("source")
			eventType, _ := cmd.Flags().GetString("event-type")
			condition, _ := cmd.Flags().GetString("condition")
			evaluationType, _ := cmd.Flags().GetString("evaluation-type")
			alertLevel, _ := cmd.Flags().GetString("alert-level")

			baseURL, _ := cmd.Flags().GetString("base-url")

			payload := map[string]interface{}{
				"name":            name,
				"description":     description,
				"source":          source,
				"event_type":      eventType,
				"condition":       condition,
				"evaluation_type": evaluationType,
				"alert_level":     alertLevel,
			}

			body, err := json.Marshal(payload)
			if err != nil {
				return err
			}

			resp, err := http.Post(baseURL+"/api/rules", "application/json", bytes.NewBuffer(body))
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
				respBody, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("server returned status: %s, body: %s", resp.Status, string(respBody))
			}

			var rule map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&rule); err != nil {
				return fmt.Errorf("failed to decode response: %w", err)
			}

			// Print table with Attribute/Value columns
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

	cmd.Flags().StringP("name", "n", "", "Name of the rule")
	cmd.Flags().StringP("description", "d", "", "Description of the rule")
	cmd.Flags().StringP("source", "s", "", "Source of the rule")
	cmd.Flags().StringP("event-type", "e", "", "Event type for the rule")
	cmd.Flags().StringP("condition", "c", "", "Sql for the rule")
	cmd.Flags().StringP("evaluation-type", "t", "", "Evaluation type for the rule")
	cmd.Flags().StringP("alert-level", "a", "", "Alert level for the rule")
	cmd.Flags().BoolP("active", "A", true, "Whether the rule is active")

	return cmd
}
