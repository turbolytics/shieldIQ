package github

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
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
