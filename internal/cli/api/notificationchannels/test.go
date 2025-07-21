package notificationchannels

import (
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

func NewTestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test <notification-channel-id>",
		Short: "Test a notification channel by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			baseURL, _ := cmd.Flags().GetString("base-url")
			url := fmt.Sprintf("%s/api/notification-channels/%s/test", baseURL, id)
			resp, err := http.Post(url, "application/json", nil)
			if err != nil {
				return fmt.Errorf("failed to call test endpoint: %w", err)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read response body: %w", err)
			}
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("test failed: %s", string(body))
			}
			fmt.Println("Test successful:", string(body))
			return nil
		},
	}
	return cmd
}
