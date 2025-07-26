package rules

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
	"github.com/turbolytics/sqlsec/internal"
)

// Pre-made rules map: name -> SQL
var preMadeRules = map[string]internal.Rule{
	"github-pull-request-merged-no-reviewers": {
		Name:           "no-reviewers",
		Description:    "Detects pull requests that were merged without any reviewers",
		EvaluationType: internal.EvaluationTypeLiveTrigger,
		EventSource:    "github",
		EventType:      "github.pull_request",
		SQL: `
SELECT *
FROM events
WHERE
  raw_payload->>'action' = 'closed'
  AND (raw_payload->'pull_request'->>'merged')::boolean = true
  AND jsonb_array_length(raw_payload->'pull_request'->'assignees') = 0
  AND jsonb_array_length(raw_payload->'pull_request'->'requested_reviewers') = 0
  AND (raw_payload->'pull_request'->>'comments')::int = 0;
`,
		AlertLevel: internal.AlertLevelLow,
		Active:     true,
	},
}

func NewInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install a pre-made rule",
		RunE: func(cmd *cobra.Command, args []string) error {
			list, _ := cmd.Flags().GetBool("list")
			if list {
				// Print table of all preMadeRules using go-pretty, following the create.go example
				t := table.NewWriter()
				t.SetOutputMirror(cmd.OutOrStdout())
				t.AppendHeader(table.Row{"ID", "Name", "Description", "Event Source", "Event Type", "Evaluation Type", "Alert Level", "Active"})
				for id, rule := range preMadeRules {
					t.AppendRow(table.Row{
						id,
						rule.Name,
						rule.Description,
						rule.EventSource,
						rule.EventType,
						rule.EvaluationType,
						rule.AlertLevel,
						rule.Active,
					})
				}
				t.SetStyle(table.StyleDefault)
				t.Style().Options.SeparateRows = false
				t.Style().Box = table.StyleBoxDefault
				t.Style().Format.Header = text.FormatDefault
				t.Render()
				return nil
			}
			id, _ := cmd.Flags().GetString("id")
			rule, ok := preMadeRules[id]
			if !ok {
				return fmt.Errorf("Unknown rule: %s", id)
			}
			// Call the create logic: reuse the create command's RunE
			baseURL, _ := cmd.Flags().GetString("base-url")
			// Build args for create command instead of setting flags directly
			argsList := []string{
				"--name", rule.Name,
				"--description", rule.Description,
				"--condition", rule.SQL,
				"--source", rule.EventSource,
				"--event-type", rule.EventType,
				"--evaluation-type", string(rule.EvaluationType),
				"--alert-level", string(rule.AlertLevel),
			}
			createCmd := NewCreateCmd()
			createCmd.SetArgs(argsList)
			if baseURL != "" {
				createCmd.Flags().Set("base-url", baseURL)
			}
			return createCmd.Execute()
		},
	}
	cmd.Flags().String("id", "", "Id of the pre-made rule")
	cmd.Flags().Bool("list", false, "List all available pre-made rules")
	cmd.MarkFlagsMutuallyExclusive("id", "list")
	return cmd
}
