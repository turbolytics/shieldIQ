package rules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

func NewUpdateCmd() *cobra.Command {
	var active bool
	cmd := &cobra.Command{
		Use:   "update <rule-id>",
		Short: "Update a rule's active status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ruleID := args[0]
			baseURL, _ := cmd.Flags().GetString("base-url")

			payload := map[string]interface{}{
				"active": active,
			}
			body, err := json.Marshal(payload)
			if err != nil {
				return err
			}

			url := fmt.Sprintf("%s/api/rules/%s", baseURL, ruleID)
			req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(body))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("server returned status: %s", resp.Status)
			}

			fmt.Println("Rule updated successfully.")
			return nil
		},
	}
	cmd.Flags().BoolVar(&active, "active", false, "Set rule active status (true/false)")
	return cmd
}
