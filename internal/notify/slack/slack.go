package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/turbolytics/sqlsec/internal/notify"
)

type SlackNotifier struct{}

func (s *SlackNotifier) Test(ctx context.Context, cfg map[string]string) error {
	webhook, ok := cfg["webhook_url"]
	if !ok || webhook == "" {
		return errors.New("missing webhook_url in config")
	}
	msg := notify.Message{
		Title: "Test Notification",
		Body:  "This is a test notification from SQLSec.",
	}
	return s.Send(ctx, cfg, msg)
}

var slackMessageTemplate = "*{{.Title}}* _(Severity: {{.RuleAlertLevel | upper}})_\n\n" +
	"{{.Body}}\n\n" +
	">*Event Details:*\n" +
	"• *Source:* {{.EventSource}}\n" +
	"• *Type:* {{.EventType}}\n" +
	"• *Resource:* <{{.ResourceLink}}|View Resource>\n\n" +
	">*Rule Triggered:*\n" +
	"• *Name:* {{.RuleName}}\n" +
	"• *Description:* {{.RuleDescription}}\n" +
	"• *Evaluation Type:* {{.RuleEvaluationType}}\n\n" +
	">*Rule Logic (SQL):*\n" +
	"```\n" +
	"{{.RuleSQL}}\n" +
	"```"

func (s *SlackNotifier) Send(ctx context.Context, cfg map[string]string, msg notify.Message) error {
	webhook, ok := cfg["webhook_url"]
	if !ok || webhook == "" {
		return errors.New("missing webhook_url in config")
	}

	tmpl, err := template.New("slackMessage").
		Funcs(template.FuncMap{
			"upper": strings.ToUpper,
		}).
		Parse(slackMessageTemplate)

	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, msg); err != nil {
		return err
	}

	payload := map[string]string{
		"text": buf.String(),
	}
	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", webhook, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("slack webhook returned status: " + resp.Status)
	}
	return nil
}

func InitializeSlack(reg *notify.Registry) {
	reg.Register(notify.SlackChannel, &SlackNotifier{})
}
