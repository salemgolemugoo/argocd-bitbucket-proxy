package forwarder

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestForward_SendsCorrectHeadersAndBody(t *testing.T) {
	var receivedHeaders http.Header
	var receivedBody []byte

	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeaders = r.Header
		receivedBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer mock.Close()

	fwd := New(mock.URL)

	payload := map[string]string{"ref": "refs/heads/main"}
	body, _ := json.Marshal(payload)

	err := fwd.Forward("push", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedHeaders.Get("X-GitHub-Event") != "push" {
		t.Errorf("expected X-GitHub-Event=push, got '%s'", receivedHeaders.Get("X-GitHub-Event"))
	}
	if receivedHeaders.Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type=application/json, got '%s'", receivedHeaders.Get("Content-Type"))
	}

	var got map[string]string
	if err := json.Unmarshal(receivedBody, &got); err != nil {
		t.Fatalf("failed to parse received body: %v", err)
	}
	if got["ref"] != "refs/heads/main" {
		t.Errorf("expected ref=refs/heads/main, got '%s'", got["ref"])
	}
}

func TestForward_ServerError(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mock.Close()

	fwd := New(mock.URL)
	err := fwd.Forward("push", []byte(`{}`))
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestForward_ConnectionError(t *testing.T) {
	fwd := New("http://localhost:1")
	err := fwd.Forward("push", []byte(`{}`))
	if err == nil {
		t.Fatal("expected error for connection failure")
	}
}
