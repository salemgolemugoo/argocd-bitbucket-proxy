package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/salemgolemugoo/argocd-bitbucket-proxy/internal/config"
)

func sign(body []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func TestHealthz(t *testing.T) {
	cfg := &config.Config{BitbucketServerSecret: "s"}
	srv := New(cfg, "http://localhost:9999")
	handler := srv.Handler()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestReadyz_WithSecret(t *testing.T) {
	cfg := &config.Config{BitbucketServerSecret: "s"}
	srv := New(cfg, "http://localhost:9999")
	handler := srv.Handler()

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestWebhook_MissingEventKey(t *testing.T) {
	cfg := &config.Config{BitbucketServerSecret: "s"}
	srv := New(cfg, "http://localhost:9999")
	handler := srv.Handler()

	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(`{}`))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestWebhook_TrailingSlash(t *testing.T) {
	cfg := &config.Config{BitbucketServerSecret: "s"}
	srv := New(cfg, "http://localhost:9999")
	handler := srv.Handler()

	req := httptest.NewRequest(http.MethodPost, "/webhook/", strings.NewReader(`{}`))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code == http.StatusNotFound {
		t.Errorf("POST /webhook/ should not return 404, got %d", w.Code)
	}
}

func TestWebhook_InvalidSignature(t *testing.T) {
	cfg := &config.Config{BitbucketServerSecret: "correct-secret"}
	srv := New(cfg, "http://localhost:9999")
	handler := srv.Handler()

	body := `{"eventKey":"repo:refs_changed","repository":{"slug":"r","project":{"key":"P"}},"changes":[]}`
	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(body))
	req.Header.Set("X-Event-Key", "repo:refs_changed")
	req.Header.Set("X-Hub-Signature", "sha256=wrong")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestWebhook_ServerPush_EndToEnd(t *testing.T) {
	var receivedEvent string
	var receivedBody []byte

	mockArgoCD := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedEvent = r.Header.Get("X-GitHub-Event")
		receivedBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer mockArgoCD.Close()

	cfg := &config.Config{BitbucketServerSecret: "test-secret"}
	srv := New(cfg, mockArgoCD.URL)
	handler := srv.Handler()

	body := []byte(`{
		"eventKey": "repo:refs_changed",
		"repository": {
			"slug": "my-repo", "name": "My Repo",
			"project": {"key": "PROJ", "name": "Project"},
			"links": {
				"clone": [{"href": "https://bb.example.com/scm/proj/my-repo.git", "name": "http"}, {"href": "ssh://git@bb.example.com:7999/proj/my-repo.git", "name": "ssh"}],
				"self": [{"href": "https://bb.example.com/projects/PROJ/repos/my-repo/browse"}]
			}
		},
		"changes": [{"ref": {"id": "refs/heads/main", "displayId": "main", "type": "BRANCH"}, "refId": "refs/heads/main", "fromHash": "aaa", "toHash": "bbb", "type": "UPDATE"}]
	}`)

	sig := sign(body, "test-secret")

	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(string(body)))
	req.Header.Set("X-Event-Key", "repo:refs_changed")
	req.Header.Set("X-Hub-Signature", sig)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if receivedEvent != "push" {
		t.Errorf("expected X-GitHub-Event=push, got '%s'", receivedEvent)
	}

	var ghPayload map[string]interface{}
	if err := json.Unmarshal(receivedBody, &ghPayload); err != nil {
		t.Fatalf("failed to parse forwarded body: %v", err)
	}
	if ghPayload["ref"] != "refs/heads/main" {
		t.Errorf("expected ref=refs/heads/main, got '%v'", ghPayload["ref"])
	}
}

func TestWebhook_CloudPush_EndToEnd(t *testing.T) {
	var receivedEvent string

	mockArgoCD := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedEvent = r.Header.Get("X-GitHub-Event")
		w.WriteHeader(http.StatusOK)
	}))
	defer mockArgoCD.Close()

	cfg := &config.Config{BitbucketCloudSecret: "cloud-secret"}
	srv := New(cfg, mockArgoCD.URL)
	handler := srv.Handler()

	body := []byte(`{
		"repository": {
			"full_name": "ws/repo", "name": "repo",
			"links": {"html": {"href": "https://bitbucket.org/ws/repo"}},
			"owner": {"nickname": "ws"}
		},
		"push": {"changes": [{"new": {"type": "branch", "name": "main", "target": {"hash": "abc"}}, "old": {"type": "branch", "name": "main", "target": {"hash": "def"}}}]}
	}`)

	sig := sign(body, "cloud-secret")

	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(string(body)))
	req.Header.Set("X-Event-Key", "repo:push")
	req.Header.Set("X-Hub-Signature", sig)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if receivedEvent != "push" {
		t.Errorf("expected X-GitHub-Event=push, got '%s'", receivedEvent)
	}
}
