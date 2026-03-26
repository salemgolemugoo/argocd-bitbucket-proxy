package translator

import (
	"os"
	"testing"
)

func TestTranslateBitbucketCloudPush(t *testing.T) {
	body, err := os.ReadFile("../../testdata/bitbucket_cloud_push.json")
	if err != nil {
		t.Fatalf("failed to read test fixture: %v", err)
	}

	event, payload, err := TranslateBitbucketCloudPush(body)
	if err != nil {
		t.Fatalf("translation error: %v", err)
	}

	if event != "push" {
		t.Errorf("expected event 'push', got '%s'", event)
	}

	push, ok := payload.(*GitHubPushPayload)
	if !ok {
		t.Fatalf("expected *GitHubPushPayload, got %T", payload)
	}

	if push.Ref != "refs/heads/main" {
		t.Errorf("expected ref 'refs/heads/main', got '%s'", push.Ref)
	}
	if push.After != "abc123def456" {
		t.Errorf("expected after 'abc123def456', got '%s'", push.After)
	}
	if push.Before != "000111222333" {
		t.Errorf("expected before '000111222333', got '%s'", push.Before)
	}
	if push.Repository.FullName != "myworkspace/my-repo" {
		t.Errorf("expected full_name 'myworkspace/my-repo', got '%s'", push.Repository.FullName)
	}
	if push.Repository.Owner.Login != "myworkspace" {
		t.Errorf("expected owner 'myworkspace', got '%s'", push.Repository.Owner.Login)
	}
	if push.Repository.CloneURL != "https://bitbucket.org/myworkspace/my-repo.git" {
		t.Errorf("expected clone_url with .git suffix, got '%s'", push.Repository.CloneURL)
	}
	if push.Repository.SSHURL != "git@bitbucket.org:myworkspace/my-repo.git" {
		t.Errorf("expected ssh_url 'git@bitbucket.org:myworkspace/my-repo.git', got '%s'", push.Repository.SSHURL)
	}
}

func TestTranslateBitbucketCloudPush_BranchCreated(t *testing.T) {
	body := []byte(`{
		"repository": {
			"full_name": "ws/repo", "name": "repo",
			"links": {"html": {"href": "https://bitbucket.org/ws/repo"}},
			"owner": {"nickname": "ws"}
		},
		"push": {
			"changes": [{
				"new": {"type": "branch", "name": "feature", "target": {"hash": "newcommit123"}},
				"old": null
			}]
		}
	}`)

	event, payload, err := TranslateBitbucketCloudPush(body)
	if err != nil {
		t.Fatalf("translation error: %v", err)
	}
	if event != "push" {
		t.Errorf("expected 'push', got '%s'", event)
	}
	push := payload.(*GitHubPushPayload)
	if push.Before != "0000000000000000000000000000000000000000" {
		t.Errorf("expected zero hash for new branch, got '%s'", push.Before)
	}
	if push.Ref != "refs/heads/feature" {
		t.Errorf("expected ref 'refs/heads/feature', got '%s'", push.Ref)
	}
}

func TestTranslateBitbucketCloudPush_NoChanges(t *testing.T) {
	body := []byte(`{
		"repository": {"full_name": "ws/repo", "name": "repo", "links": {"html": {"href": "https://bitbucket.org/ws/repo"}}, "owner": {"nickname": "ws"}},
		"push": {"changes": []}
	}`)
	_, _, err := TranslateBitbucketCloudPush(body)
	if err == nil {
		t.Fatal("expected error for empty changes")
	}
}

func TestIsBitbucketCloud(t *testing.T) {
	cloudBody := []byte(`{"repository": {"full_name": "ws/repo"}, "push": {"changes": []}}`)
	if !IsBitbucketCloud(cloudBody) {
		t.Error("expected IsBitbucketCloud=true for cloud payload")
	}

	serverBody := []byte(`{"eventKey": "repo:refs_changed", "repository": {"slug": "repo", "project": {"key": "PROJ"}}}`)
	if IsBitbucketCloud(serverBody) {
		t.Error("expected IsBitbucketCloud=false for server payload")
	}

	emptyBody := []byte(`{}`)
	if IsBitbucketCloud(emptyBody) {
		t.Error("expected IsBitbucketCloud=false for empty payload")
	}
}

func TestTranslateBitbucketCloudPR_Created(t *testing.T) {
	body, err := os.ReadFile("../../testdata/bitbucket_cloud_pr_created.json")
	if err != nil {
		t.Fatalf("failed to read test fixture: %v", err)
	}

	event, payload, err := TranslateBitbucketCloudPR(body, "pullrequest:created")
	if err != nil {
		t.Fatalf("translation error: %v", err)
	}

	if event != "pull_request" {
		t.Errorf("expected event 'pull_request', got '%s'", event)
	}

	pr, ok := payload.(*GitHubPullRequestPayload)
	if !ok {
		t.Fatalf("expected *GitHubPullRequestPayload, got %T", payload)
	}

	if pr.Action != "opened" {
		t.Errorf("expected action 'opened', got '%s'", pr.Action)
	}
	if pr.Number != 7 {
		t.Errorf("expected number 7, got %d", pr.Number)
	}
	if pr.PullRequest.Head.Ref != "feature/cool" {
		t.Errorf("expected head ref 'feature/cool', got '%s'", pr.PullRequest.Head.Ref)
	}
	if pr.PullRequest.Base.Ref != "main" {
		t.Errorf("expected base ref 'main', got '%s'", pr.PullRequest.Base.Ref)
	}
	if pr.PullRequest.Head.Repo.CloneURL != "https://bitbucket.org/myworkspace/my-repo.git" {
		t.Errorf("expected head clone_url with .git suffix, got '%s'", pr.PullRequest.Head.Repo.CloneURL)
	}
	if pr.PullRequest.Base.Repo.CloneURL != "https://bitbucket.org/myworkspace/my-repo.git" {
		t.Errorf("expected base clone_url with .git suffix, got '%s'", pr.PullRequest.Base.Repo.CloneURL)
	}
}

func TestTranslateBitbucketCloudPR_Fulfilled(t *testing.T) {
	body, err := os.ReadFile("../../testdata/bitbucket_cloud_pr_created.json")
	if err != nil {
		t.Fatalf("failed to read test fixture: %v", err)
	}

	_, payload, err := TranslateBitbucketCloudPR(body, "pullrequest:fulfilled")
	if err != nil {
		t.Fatalf("translation error: %v", err)
	}

	pr := payload.(*GitHubPullRequestPayload)
	if pr.Action != "closed" {
		t.Errorf("expected action 'closed', got '%s'", pr.Action)
	}
	if !pr.PullRequest.Merged {
		t.Errorf("expected merged=true")
	}
}

func TestTranslateBitbucketCloudPR_Rejected(t *testing.T) {
	body, err := os.ReadFile("../../testdata/bitbucket_cloud_pr_created.json")
	if err != nil {
		t.Fatalf("failed to read test fixture: %v", err)
	}

	_, payload, err := TranslateBitbucketCloudPR(body, "pullrequest:rejected")
	if err != nil {
		t.Fatalf("translation error: %v", err)
	}

	pr := payload.(*GitHubPullRequestPayload)
	if pr.Action != "closed" {
		t.Errorf("expected action 'closed', got '%s'", pr.Action)
	}
	if pr.PullRequest.Merged {
		t.Errorf("expected merged=false")
	}
}
