package translator

import (
	"os"
	"testing"
)

func TestTranslate_BitbucketServerPush(t *testing.T) {
	body, err := os.ReadFile("../../testdata/bitbucket_server_push.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}
	event, payload, err := Translate(body, "repo:refs_changed")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event != "push" {
		t.Errorf("expected 'push', got '%s'", event)
	}
	if _, ok := payload.(*GitHubPushPayload); !ok {
		t.Errorf("expected *GitHubPushPayload, got %T", payload)
	}
}

func TestTranslate_BitbucketServerPR(t *testing.T) {
	body, err := os.ReadFile("../../testdata/bitbucket_server_pr_opened.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}
	event, payload, err := Translate(body, "pr:opened")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event != "pull_request" {
		t.Errorf("expected 'pull_request', got '%s'", event)
	}
	if _, ok := payload.(*GitHubPullRequestPayload); !ok {
		t.Errorf("expected *GitHubPullRequestPayload, got %T", payload)
	}
}

func TestTranslate_BitbucketCloudPush(t *testing.T) {
	body, err := os.ReadFile("../../testdata/bitbucket_cloud_push.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}
	event, payload, err := Translate(body, "repo:push")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event != "push" {
		t.Errorf("expected 'push', got '%s'", event)
	}
	if _, ok := payload.(*GitHubPushPayload); !ok {
		t.Errorf("expected *GitHubPushPayload, got %T", payload)
	}
}

func TestTranslate_BitbucketCloudPR(t *testing.T) {
	body, err := os.ReadFile("../../testdata/bitbucket_cloud_pr_created.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}
	event, payload, err := Translate(body, "pullrequest:created")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if event != "pull_request" {
		t.Errorf("expected 'pull_request', got '%s'", event)
	}
	if _, ok := payload.(*GitHubPullRequestPayload); !ok {
		t.Errorf("expected *GitHubPullRequestPayload, got %T", payload)
	}
}

func TestTranslate_UnsupportedEvent(t *testing.T) {
	body := []byte(`{}`)
	_, _, err := Translate(body, "repo:comment:added")
	if err == nil {
		t.Fatal("expected error for unsupported event")
	}
}
