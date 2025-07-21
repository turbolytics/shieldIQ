package notificationchannels

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all notification channels",
		RunE: func(cmd *cobra.Command, args []string) error {
			baseURL, _ := cmd.Flags().GetString("base-url")
			resp, err := http.Get(baseURL + "/api/notification-channels")
			if err != nil {
				return fmt.Errorf("failed to call list endpoint: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("server returned status: %s", resp.Status)
			}
			var channels []map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&channels); err != nil {
				return fmt.Errorf("failed to decode response: %w", err)
			}
			if len(channels) == 0 {
				fmt.Println("No notification channels found.")
				return nil
			}
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			// Collect all unique keys for columns
			columns := []string{"id", "name", "type", "config", "created_at", "updated_at"}
			header := make(table.Row, len(columns))
			for i, col := range columns {
				header[i] = col
			}
			t.AppendHeader(header)
			for _, ch := range channels {
				row := make(table.Row, len(columns))
				for i, col := range columns {
					if val, ok := ch[col]; ok {
						if col == "config" {
							b, _ := json.Marshal(val)
							row[i] = string(b[:50])
						} else {
							row[i] = val
						}
					} else {
						row[i] = ""
					}
				}
				t.AppendRow(row)
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
