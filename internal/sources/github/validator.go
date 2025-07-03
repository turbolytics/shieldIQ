package github

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
)

type GithubValidator struct{}

func (v *GithubValidator) Validate(r *http.Request, secret string) error {
	sig := r.Header.Get("X-Hub-Signature-256")
	if sig == "" {
		return errors.New("missing signature header")
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	// Restore body for downstream use
	r.Body = io.NopCloser(io.MultiReader(bytes.NewReader(body)))

	h := hmac.New(sha256.New, []byte(secret))
	h.Write(body)
	expected := "sha256=" + hex.EncodeToString(h.Sum(nil))
	if !hmac.Equal([]byte(sig), []byte(expected)) {
		return errors.New("invalid signature")
	}
	return nil
}
