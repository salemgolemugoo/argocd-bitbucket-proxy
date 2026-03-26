package translator

import (
	"encoding/json"
	"fmt"
	"strings"
)

func IsBitbucketCloud(body []byte) bool {
	var probe struct {
		Repository struct {
			FullName string `json:"full_name"`
		} `json:"repository"`
		EventKey string `json:"eventKey"`
	}
	if err := json.Unmarshal(body, &probe); err != nil {
		return false
	}
	return probe.Repository.FullName != "" && probe.EventKey == ""
}

func TranslateBitbucketCloudPush(body []byte) (string, interface{}, error) {
	var payload BitbucketCloudPushPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", nil, fmt.Errorf("failed to parse Bitbucket Cloud push payload: %w", err)
	}

	if len(payload.Push.Changes) == 0 {
		return "", nil, fmt.Errorf("no changes in push payload")
	}

	change := payload.Push.Changes[0]
	if change.New == nil {
		return "", nil, fmt.Errorf("change has no new ref (branch deleted)")
	}

	var before string
	if change.Old != nil {
		before = change.Old.Target.Hash
	} else {
		before = "0000000000000000000000000000000000000000"
	}

	htmlURL := payload.Repository.Links.HTML.Href
	cloneURL := htmlURL + ".git"
	sshURL := buildCloudSSHURL(payload.Repository.FullName)
	owner, name := splitFullName(payload.Repository.FullName)

	ghPayload := &GitHubPushPayload{
		Ref:    "refs/heads/" + change.New.Name,
		Before: before,
		After:  change.New.Target.Hash,
		Repository: GitHubRepository{
			Name:     name,
			FullName: payload.Repository.FullName,
			Owner:    GitHubUser{Login: owner},
			HTMLURL:  htmlURL,
			CloneURL: cloneURL,
			SSHURL:   sshURL,
		},
	}

	return "push", ghPayload, nil
}

func TranslateBitbucketCloudPR(body []byte, eventKey string) (string, interface{}, error) {
	var payload BitbucketCloudPRPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", nil, fmt.Errorf("failed to parse Bitbucket Cloud PR payload: %w", err)
	}

	pr := payload.PullRequest
	action, merged := mapCloudPREvent(eventKey)

	srcHTMLURL := pr.Source.Repository.Links.HTML.Href
	srcCloneURL := srcHTMLURL + ".git"
	srcSSHURL := buildCloudSSHURL(pr.Source.Repository.FullName)
	srcOwner, srcName := splitFullName(pr.Source.Repository.FullName)

	dstHTMLURL := pr.Destination.Repository.Links.HTML.Href
	dstCloneURL := dstHTMLURL + ".git"
	dstSSHURL := buildCloudSSHURL(pr.Destination.Repository.FullName)
	dstOwner, dstName := splitFullName(pr.Destination.Repository.FullName)

	ghPayload := &GitHubPullRequestPayload{
		Action: action,
		Number: pr.ID,
		PullRequest: GitHubPullRequest{
			Title:  pr.Title,
			State:  mapCloudPRStateToGitHub(pr.State),
			Merged: merged,
			Head: GitHubPRBranch{
				Ref: pr.Source.Branch.Name,
				SHA: pr.Source.Commit.Hash,
				Repo: GitHubRepository{
					Name:     srcName,
					FullName: pr.Source.Repository.FullName,
					Owner:    GitHubUser{Login: srcOwner},
					HTMLURL:  srcHTMLURL,
					CloneURL: srcCloneURL,
					SSHURL:   srcSSHURL,
				},
			},
			Base: GitHubPRBranch{
				Ref: pr.Destination.Branch.Name,
				SHA: pr.Destination.Commit.Hash,
				Repo: GitHubRepository{
					Name:     dstName,
					FullName: pr.Destination.Repository.FullName,
					Owner:    GitHubUser{Login: dstOwner},
					HTMLURL:  dstHTMLURL,
					CloneURL: dstCloneURL,
					SSHURL:   dstSSHURL,
				},
			},
		},
		Repository: GitHubRepository{
			Name:     dstName,
			FullName: pr.Destination.Repository.FullName,
			Owner:    GitHubUser{Login: dstOwner},
			HTMLURL:  dstHTMLURL,
			CloneURL: dstCloneURL,
			SSHURL:   dstSSHURL,
		},
	}

	return "pull_request", ghPayload, nil
}

func splitFullName(fullName string) (owner, name string) {
	parts := strings.SplitN(fullName, "/", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", fullName
}

func buildCloudSSHURL(fullName string) string {
	return "git@bitbucket.org:" + fullName + ".git"
}

func mapCloudPREvent(eventKey string) (action string, merged bool) {
	switch eventKey {
	case "pullrequest:fulfilled":
		return "closed", true
	case "pullrequest:rejected":
		return "closed", false
	case "pullrequest:updated":
		return "synchronize", false
	default:
		return "opened", false
	}
}

func mapCloudPRStateToGitHub(state string) string {
	switch strings.ToUpper(state) {
	case "OPEN":
		return "open"
	default:
		return "closed"
	}
}
