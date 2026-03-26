package translator

import (
	"encoding/json"
	"fmt"
)

func IsBitbucketServer(body []byte) bool {
	var probe struct {
		EventKey   string `json:"eventKey"`
		Repository struct {
			Slug    string `json:"slug"`
			Project struct {
				Key string `json:"key"`
			} `json:"project"`
		} `json:"repository"`
	}
	if err := json.Unmarshal(body, &probe); err != nil {
		return false
	}
	return probe.EventKey != "" && probe.Repository.Slug != "" && probe.Repository.Project.Key != ""
}

func TranslateBitbucketServerPush(body []byte) (string, interface{}, error) {
	var payload BitbucketServerPushPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", nil, fmt.Errorf("failed to parse Bitbucket Server push payload: %w", err)
	}

	if len(payload.Changes) == 0 {
		return "", nil, fmt.Errorf("no changes in push payload")
	}

	change := payload.Changes[0]
	cloneURL, sshURL := extractServerCloneURLs(payload.Repository.Links)
	selfURL := extractServerSelfURL(payload.Repository.Links)

	ghPayload := &GitHubPushPayload{
		Ref:    change.RefID,
		Before: change.FromHash,
		After:  change.ToHash,
		Repository: GitHubRepository{
			Name:     payload.Repository.Slug,
			FullName: payload.Repository.Project.Key + "/" + payload.Repository.Slug,
			Owner:    GitHubUser{Login: payload.Repository.Project.Key},
			HTMLURL:  selfURL,
			CloneURL: cloneURL,
			SSHURL:   sshURL,
		},
	}

	return "push", ghPayload, nil
}

func TranslateBitbucketServerPR(body []byte) (string, interface{}, error) {
	var payload BitbucketServerPRPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", nil, fmt.Errorf("failed to parse Bitbucket Server PR payload: %w", err)
	}

	pr := payload.PullRequest
	action, merged := mapServerPRState(pr.State, pr.Open)
	fromCloneURL, fromSSHURL := extractServerCloneURLs(pr.FromRef.Repository.Links)
	fromSelfURL := extractServerSelfURL(pr.FromRef.Repository.Links)
	toCloneURL, toSSHURL := extractServerCloneURLs(pr.ToRef.Repository.Links)
	toSelfURL := extractServerSelfURL(pr.ToRef.Repository.Links)

	ghPayload := &GitHubPullRequestPayload{
		Action: action,
		Number: int64(pr.ID),
		PullRequest: GitHubPullRequest{
			Title:  pr.Title,
			State:  mapPRStateToGitHub(pr.State),
			Merged: merged,
			Head: GitHubPRBranch{
				Ref: pr.FromRef.DisplayID,
				Repo: GitHubRepository{
					Name:     pr.FromRef.Repository.Slug,
					FullName: pr.FromRef.Repository.Project.Key + "/" + pr.FromRef.Repository.Slug,
					Owner:    GitHubUser{Login: pr.FromRef.Repository.Project.Key},
					HTMLURL:  fromSelfURL,
					CloneURL: fromCloneURL,
					SSHURL:   fromSSHURL,
				},
			},
			Base: GitHubPRBranch{
				Ref: pr.ToRef.DisplayID,
				Repo: GitHubRepository{
					Name:     pr.ToRef.Repository.Slug,
					FullName: pr.ToRef.Repository.Project.Key + "/" + pr.ToRef.Repository.Slug,
					Owner:    GitHubUser{Login: pr.ToRef.Repository.Project.Key},
					HTMLURL:  toSelfURL,
					CloneURL: toCloneURL,
					SSHURL:   toSSHURL,
				},
			},
		},
		Repository: GitHubRepository{
			Name:     pr.ToRef.Repository.Slug,
			FullName: pr.ToRef.Repository.Project.Key + "/" + pr.ToRef.Repository.Slug,
			Owner:    GitHubUser{Login: pr.ToRef.Repository.Project.Key},
			HTMLURL:  toSelfURL,
			CloneURL: toCloneURL,
			SSHURL:   toSSHURL,
		},
	}

	return "pull_request", ghPayload, nil
}

func extractServerCloneURLs(links map[string]interface{}) (httpURL, sshURL string) {
	cloneRaw, ok := links["clone"]
	if !ok {
		return "", ""
	}
	clones, ok := cloneRaw.([]interface{})
	if !ok {
		return "", ""
	}
	for _, c := range clones {
		m, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		name, _ := m["name"].(string)
		href, _ := m["href"].(string)
		switch name {
		case "http":
			httpURL = href
		case "ssh":
			sshURL = href
		}
	}
	return httpURL, sshURL
}

func extractServerSelfURL(links map[string]interface{}) string {
	selfRaw, ok := links["self"]
	if !ok {
		return ""
	}
	selfs, ok := selfRaw.([]interface{})
	if !ok {
		return ""
	}
	if len(selfs) == 0 {
		return ""
	}
	m, ok := selfs[0].(map[string]interface{})
	if !ok {
		return ""
	}
	href, _ := m["href"].(string)
	return href
}

func mapServerPRState(state string, open bool) (action string, merged bool) {
	switch state {
	case "MERGED":
		return "closed", true
	case "DECLINED":
		return "closed", false
	default:
		if open {
			return "opened", false
		}
		return "closed", false
	}
}

func mapPRStateToGitHub(state string) string {
	switch state {
	case "OPEN":
		return "open"
	default:
		return "closed"
	}
}
