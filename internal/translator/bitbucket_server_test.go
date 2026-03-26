package translator

import (
	"os"
	"testing"
)

func TestTranslateBitbucketServerPush(t *testing.T) {
	body, err := os.ReadFile("../../testdata/bitbucket_server_push.json")
	if err != nil {
		t.Fatalf("failed to read test fixture: %v", err)
	}

	event, payload, err := TranslateBitbucketServerPush(body)
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
	if push.After != "eeff0011" {
		t.Errorf("expected after 'eeff0011', got '%s'", push.After)
	}
	if push.Before != "aabbccdd" {
		t.Errorf("expected before 'aabbccdd', got '%s'", push.Before)
	}
	if push.Repository.FullName != "PROJ/my-repo" {
		t.Errorf("expected full_name 'PROJ/my-repo', got '%s'", push.Repository.FullName)
	}
	if push.Repository.Name != "my-repo" {
		t.Errorf("expected name 'my-repo', got '%s'", push.Repository.Name)
	}
	if push.Repository.Owner.Login != "PROJ" {
		t.Errorf("expected owner 'PROJ', got '%s'", push.Repository.Owner.Login)
	}
	if push.Repository.CloneURL != "https://bitbucket.example.com/scm/proj/my-repo.git" {
		t.Errorf("expected clone_url, got '%s'", push.Repository.CloneURL)
	}
	if push.Repository.SSHURL != "ssh://git@bitbucket.example.com:7999/proj/my-repo.git" {
		t.Errorf("expected ssh_url, got '%s'", push.Repository.SSHURL)
	}
	if push.Repository.HTMLURL != "https://bitbucket.example.com/projects/PROJ/repos/my-repo/browse" {
		t.Errorf("expected html_url, got '%s'", push.Repository.HTMLURL)
	}
}

func TestTranslateBitbucketServerPush_MultipleChanges(t *testing.T) {
	body := []byte(`{
		"eventKey": "repo:refs_changed",
		"repository": {
			"slug": "repo", "name": "repo",
			"project": {"key": "PROJ", "name": "Project"},
			"links": {"clone": [{"href": "https://bb.example.com/scm/proj/repo.git", "name": "http"}, {"href": "ssh://git@bb.example.com:7999/proj/repo.git", "name": "ssh"}], "self": [{"href": "https://bb.example.com/projects/PROJ/repos/repo/browse"}]}
		},
		"changes": [
			{"ref": {"id": "refs/heads/main", "displayId": "main", "type": "BRANCH"}, "refId": "refs/heads/main", "fromHash": "aaaa", "toHash": "bbbb", "type": "UPDATE"},
			{"ref": {"id": "refs/heads/develop", "displayId": "develop", "type": "BRANCH"}, "refId": "refs/heads/develop", "fromHash": "cccc", "toHash": "dddd", "type": "UPDATE"}
		]
	}`)

	event, payload, err := TranslateBitbucketServerPush(body)
	if err != nil {
		t.Fatalf("translation error: %v", err)
	}
	if event != "push" {
		t.Errorf("expected event 'push', got '%s'", event)
	}
	push := payload.(*GitHubPushPayload)
	if push.Ref != "refs/heads/main" {
		t.Errorf("expected ref from first change, got '%s'", push.Ref)
	}
}

func TestTranslateBitbucketServerPush_NoChanges(t *testing.T) {
	body := []byte(`{
		"eventKey": "repo:refs_changed",
		"repository": {"slug": "repo", "name": "repo", "project": {"key": "PROJ", "name": "Project"}, "links": {"clone": [], "self": []}},
		"changes": []
	}`)
	_, _, err := TranslateBitbucketServerPush(body)
	if err == nil {
		t.Fatal("expected error for empty changes")
	}
}

func TestTranslateBitbucketServerPR_Opened(t *testing.T) {
	body, err := os.ReadFile("../../testdata/bitbucket_server_pr_opened.json")
	if err != nil {
		t.Fatalf("failed to read test fixture: %v", err)
	}

	event, payload, err := TranslateBitbucketServerPR(body)
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
	if pr.Number != 42 {
		t.Errorf("expected number 42, got %d", pr.Number)
	}
	if pr.PullRequest.Title != "Add feature X" {
		t.Errorf("expected title 'Add feature X', got '%s'", pr.PullRequest.Title)
	}
	if pr.PullRequest.State != "open" {
		t.Errorf("expected state 'open', got '%s'", pr.PullRequest.State)
	}
	if pr.PullRequest.Merged {
		t.Errorf("expected merged=false")
	}
	if pr.PullRequest.Head.Ref != "feature/x" {
		t.Errorf("expected head ref 'feature/x', got '%s'", pr.PullRequest.Head.Ref)
	}
	if pr.PullRequest.Base.Ref != "main" {
		t.Errorf("expected base ref 'main', got '%s'", pr.PullRequest.Base.Ref)
	}
	if pr.Repository.FullName != "PROJ/my-repo" {
		t.Errorf("expected full_name 'PROJ/my-repo', got '%s'", pr.Repository.FullName)
	}
}

func TestTranslateBitbucketServerPR_Merged(t *testing.T) {
	body := []byte(`{
		"eventKey": "pr:merged",
		"pullRequest": {
			"id": 10, "title": "Merge feature", "state": "MERGED", "open": false, "closed": true,
			"fromRef": {
				"id": "refs/heads/feature", "displayId": "feature",
				"repository": {
					"slug": "my-repo", "name": "My Repo",
					"project": {"key": "PROJ", "name": "My Project"},
					"links": {"clone": [{"href": "https://bb.example.com/scm/proj/my-repo.git", "name": "http"}, {"href": "ssh://git@bb.example.com:7999/proj/my-repo.git", "name": "ssh"}], "self": [{"href": "https://bb.example.com/projects/PROJ/repos/my-repo/browse"}]}
				}
			},
			"toRef": {
				"id": "refs/heads/main", "displayId": "main",
				"repository": {
					"slug": "my-repo", "name": "My Repo",
					"project": {"key": "PROJ", "name": "My Project"},
					"links": {"clone": [{"href": "https://bb.example.com/scm/proj/my-repo.git", "name": "http"}, {"href": "ssh://git@bb.example.com:7999/proj/my-repo.git", "name": "ssh"}], "self": [{"href": "https://bb.example.com/projects/PROJ/repos/my-repo/browse"}]}
				}
			}
		}
	}`)

	event, payload, err := TranslateBitbucketServerPR(body)
	if err != nil {
		t.Fatalf("translation error: %v", err)
	}
	if event != "pull_request" {
		t.Errorf("expected 'pull_request', got '%s'", event)
	}
	pr := payload.(*GitHubPullRequestPayload)
	if pr.Action != "closed" {
		t.Errorf("expected action 'closed', got '%s'", pr.Action)
	}
	if !pr.PullRequest.Merged {
		t.Errorf("expected merged=true")
	}
}

func TestTranslateBitbucketServerPR_Declined(t *testing.T) {
	body := []byte(`{
		"eventKey": "pr:declined",
		"pullRequest": {
			"id": 11, "title": "Declined PR", "state": "DECLINED", "open": false, "closed": true,
			"fromRef": {
				"id": "refs/heads/feature", "displayId": "feature",
				"repository": {
					"slug": "my-repo", "name": "My Repo",
					"project": {"key": "PROJ", "name": "My Project"},
					"links": {"clone": [{"href": "https://bb.example.com/scm/proj/my-repo.git", "name": "http"}, {"href": "ssh://git@bb.example.com:7999/proj/my-repo.git", "name": "ssh"}], "self": [{"href": "https://bb.example.com/projects/PROJ/repos/my-repo/browse"}]}
				}
			},
			"toRef": {
				"id": "refs/heads/main", "displayId": "main",
				"repository": {
					"slug": "my-repo", "name": "My Repo",
					"project": {"key": "PROJ", "name": "My Project"},
					"links": {"clone": [{"href": "https://bb.example.com/scm/proj/my-repo.git", "name": "http"}, {"href": "ssh://git@bb.example.com:7999/proj/my-repo.git", "name": "ssh"}], "self": [{"href": "https://bb.example.com/projects/PROJ/repos/my-repo/browse"}]}
				}
			}
		}
	}`)

	event, payload, err := TranslateBitbucketServerPR(body)
	if err != nil {
		t.Fatalf("translation error: %v", err)
	}
	if event != "pull_request" {
		t.Errorf("expected 'pull_request', got '%s'", event)
	}
	pr := payload.(*GitHubPullRequestPayload)
	if pr.Action != "closed" {
		t.Errorf("expected action 'closed', got '%s'", pr.Action)
	}
	if pr.PullRequest.Merged {
		t.Errorf("expected merged=false")
	}
}
