package webhooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"io"
	"net/http"
)

func NewWebhookCreateCmd() *cobra.Command {
	var name string
	var source string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new webhook source",
		RunE: func(cmd *cobra.Command, args []string) error {
			baseURL, _ := cmd.Flags().GetString("base-url")

			payload := map[string]interface{}{
				"name":   name,
				"source": source,
			}
			body, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("failed to marshal payload: %w", err)
			}
			url := fmt.Sprintf("%s/api/webhooks", baseURL)
			resp, err := http.Post(url, "application/json", bytes.NewReader(body))
			if err != nil {
				return fmt.Errorf("failed to call create webhook endpoint: %w", err)
			}
			defer resp.Body.Close()
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read response body: %w", err)
			}
			if resp.StatusCode != http.StatusCreated {
				return fmt.Errorf("failed to create webhook: %s", string(respBody))
			}
			var result map[string]interface{}
			if err := json.Unmarshal(respBody, &result); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			t := table.NewWriter()
			t.SetOutputMirror(cmd.OutOrStdout())
			t.AppendHeader(table.Row{"Attribute", "Value"})
			for k, v := range result {
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

	cmd.Flags().StringVar(&name, "name", "", "Name of the webhook (required)")
	cmd.Flags().StringVar(&source, "source", "", "Source system (e.g. github, auth0)")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("source")

	return cmd
}
