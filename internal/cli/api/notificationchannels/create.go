package notificationchannels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

func NewCreateCmd() *cobra.Command {
	var name string
	var channelType string
	var config string
	var configWebhookURL string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new notification channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			baseURL, _ := cmd.Parent().Parent().PersistentFlags().GetString("base-url")
			if baseURL == "" {
				return fmt.Errorf("base-url is required, set it using --base-url or shieldIQ_API_BASE_URL environment variable")
			}

			if name == "" || channelType == "" {
				return fmt.Errorf("--name and --type are required")
			}

			var configMap map[string]interface{}
			if channelType == "slack" {
				if configWebhookURL == "" {
					return fmt.Errorf("--config-webhook-url is required for type=slack")
				}
				configMap = map[string]interface{}{"webhook_url": configWebhookURL}
			} else {
				if config == "" {
					return fmt.Errorf("--config is required for non-slack types")
				}
				if err := json.Unmarshal([]byte(config), &configMap); err != nil {
					return fmt.Errorf("invalid config JSON: %w", err)
				}
			}

			payload := map[string]interface{}{
				"name":   name,
				"type":   channelType,
				"config": configMap,
			}
			body, err := json.Marshal(payload)
			if err != nil {
				return err
			}

			resp, err := http.Post(baseURL+"/api/notification-channels", "application/json", bytes.NewBuffer(body))
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
				// Read the response body to provide more context in the error
				bs, err := io.ReadAll(resp.Body)
				if err != nil {
					return fmt.Errorf("failed to decode error response: %w", err)
				}
				return fmt.Errorf("server returned status: %s, body: %s", resp.Status, string(bs))
			}

			var result map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return err
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"Attribute", "Value"})
			t.AppendRow(table.Row{"id", result["id"]})
			t.AppendRow(table.Row{"name", result["name"]})
			t.AppendRow(table.Row{"type", result["type"]})
			t.AppendRow(table.Row{"config", result["config"]})
			t.SetStyle(table.StyleLight)
			t.Style().Format.Header = text.FormatTitle
			t.Render()

			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Name of the notification channel (required)")
	cmd.Flags().StringVar(&channelType, "type", "", "Type of the notification channel (e.g. slack) (required)")
	cmd.Flags().StringVar(&config, "config", "", "Config JSON for the channel (required for non-slack)")
	cmd.Flags().StringVar(&configWebhookURL, "config-webhook-url", "", "Webhook URL for slack notification channel (required for slack)")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("type")

	return cmd
}
