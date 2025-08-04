package github

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type GithubParser struct{}

func (p *GithubParser) Parse(r *http.Request) (map[string]any, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	// Restore body for downstream use
	r.Body = io.NopCloser(io.MultiReader(bytes.NewReader(body)))
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func (p *GithubParser) Type(r *http.Request) (string, error) {
	eventType := r.Header.Get("X-GitHub-Event")
	if eventType == "" {
		return "", io.EOF
	}
	return eventType, nil
}

func (p *GithubParser) ResourceURL(payload map[string]any) (*url.URL, error) {
	pr, ok := payload["pull_request"].(map[string]any)
	if !ok {
		return nil, io.EOF
	}
	urlStr, ok := pr["html_url"].(string)
	if !ok || urlStr == "" {
		return nil, io.EOF
	}
	return url.Parse(urlStr)
}
