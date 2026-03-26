package translator

// --- Bitbucket Server types ---

type BitbucketServerPushPayload struct {
	EventKey   string                  `json:"eventKey"`
	Repository BitbucketServerRepo     `json:"repository"`
	Changes    []BitbucketServerChange `json:"changes"`
}

type BitbucketServerRepo struct {
	Slug    string                 `json:"slug"`
	Name    string                 `json:"name"`
	Project BitbucketServerProject `json:"project"`
	Links   map[string]interface{} `json:"links"`
}

type BitbucketServerProject struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type BitbucketServerChange struct {
	Ref      BitbucketServerRef `json:"ref"`
	RefID    string             `json:"refId"`
	FromHash string             `json:"fromHash"`
	ToHash   string             `json:"toHash"`
	Type     string             `json:"type"`
}

type BitbucketServerRef struct {
	ID        string `json:"id"`
	DisplayID string `json:"displayId"`
	Type      string `json:"type"`
}

type BitbucketServerPRPayload struct {
	EventKey    string            `json:"eventKey"`
	PullRequest BitbucketServerPR `json:"pullRequest"`
}

type BitbucketServerPR struct {
	ID      uint64               `json:"id"`
	Title   string               `json:"title"`
	State   string               `json:"state"`
	Open    bool                 `json:"open"`
	Closed  bool                 `json:"closed"`
	FromRef BitbucketServerPRRef `json:"fromRef"`
	ToRef   BitbucketServerPRRef `json:"toRef"`
}

type BitbucketServerPRRef struct {
	ID         string              `json:"id"`
	DisplayID  string              `json:"displayId"`
	Repository BitbucketServerRepo `json:"repository"`
}

// --- Bitbucket Cloud types ---

type BitbucketCloudPushPayload struct {
	Repository BitbucketCloudRepo `json:"repository"`
	Push       struct {
		Changes []BitbucketCloudChange `json:"changes"`
	} `json:"push"`
}

type BitbucketCloudRepo struct {
	FullName string `json:"full_name"`
	Name     string `json:"name"`
	Links    struct {
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"links"`
	Owner struct {
		Nickname string `json:"nickname"`
	} `json:"owner"`
}

type BitbucketCloudChange struct {
	New *BitbucketCloudTarget `json:"new"`
	Old *BitbucketCloudTarget `json:"old"`
}

type BitbucketCloudTarget struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	Target struct {
		Hash string `json:"hash"`
	} `json:"target"`
}

type BitbucketCloudPRPayload struct {
	Repository  BitbucketCloudRepo `json:"repository"`
	PullRequest BitbucketCloudPR   `json:"pullrequest"`
}

type BitbucketCloudPR struct {
	ID          int64                    `json:"id"`
	Title       string                   `json:"title"`
	State       string                   `json:"state"`
	Source      BitbucketCloudPREndpoint `json:"source"`
	Destination BitbucketCloudPREndpoint `json:"destination"`
}

type BitbucketCloudPREndpoint struct {
	Branch struct {
		Name string `json:"name"`
	} `json:"branch"`
	Commit struct {
		Hash string `json:"hash"`
	} `json:"commit"`
	Repository BitbucketCloudRepo `json:"repository"`
}

// --- GitHub output types (what ArgoCD expects) ---

type GitHubPushPayload struct {
	Ref        string           `json:"ref"`
	Before     string           `json:"before"`
	After      string           `json:"after"`
	Repository GitHubRepository `json:"repository"`
	Sender     GitHubUser       `json:"sender"`
}

type GitHubRepository struct {
	ID            int64      `json:"id"`
	Name          string     `json:"name"`
	FullName      string     `json:"full_name"`
	Owner         GitHubUser `json:"owner"`
	HTMLURL       string     `json:"html_url"`
	CloneURL      string     `json:"clone_url"`
	SSHURL        string     `json:"ssh_url"`
	DefaultBranch string     `json:"default_branch"`
}

type GitHubUser struct {
	Login string `json:"login"`
}

type GitHubPullRequestPayload struct {
	Action      string            `json:"action"`
	Number      int64             `json:"number"`
	PullRequest GitHubPullRequest `json:"pull_request"`
	Repository  GitHubRepository  `json:"repository"`
}

type GitHubPullRequest struct {
	Title  string         `json:"title"`
	State  string         `json:"state"`
	Merged bool           `json:"merged"`
	Head   GitHubPRBranch `json:"head"`
	Base   GitHubPRBranch `json:"base"`
}

type GitHubPRBranch struct {
	Ref  string           `json:"ref"`
	SHA  string           `json:"sha"`
	Repo GitHubRepository `json:"repo"`
}
