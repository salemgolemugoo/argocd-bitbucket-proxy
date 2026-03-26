package forwarder

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Forwarder struct {
	targetURL string
	client    *http.Client
}

func New(targetURL string) *Forwarder {
	return &Forwarder{
		targetURL: targetURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (f *Forwarder) Forward(githubEvent string, body []byte) error {
	req, err := http.NewRequest(http.MethodPost, f.targetURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Event", githubEvent)

	resp, err := f.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to forward webhook: %w", err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("ArgoCD returned status %d", resp.StatusCode)
	}

	return nil
}
