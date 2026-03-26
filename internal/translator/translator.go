package translator

import (
	"fmt"
	"strings"
)

func Translate(body []byte, eventKey string) (string, interface{}, error) {
	if IsBitbucketServer(body) {
		return translateServer(body, eventKey)
	}
	if IsBitbucketCloud(body) {
		return translateCloud(body, eventKey)
	}
	return "", nil, fmt.Errorf("unable to detect Bitbucket source from payload")
}

func translateServer(body []byte, eventKey string) (string, interface{}, error) {
	switch {
	case eventKey == "repo:refs_changed":
		return TranslateBitbucketServerPush(body)
	case strings.HasPrefix(eventKey, "pr:"):
		return TranslateBitbucketServerPR(body)
	default:
		return "", nil, fmt.Errorf("unsupported Bitbucket Server event: %s", eventKey)
	}
}

func translateCloud(body []byte, eventKey string) (string, interface{}, error) {
	switch {
	case eventKey == "repo:push":
		return TranslateBitbucketCloudPush(body)
	case strings.HasPrefix(eventKey, "pullrequest:"):
		return TranslateBitbucketCloudPR(body, eventKey)
	default:
		return "", nil, fmt.Errorf("unsupported Bitbucket Cloud event: %s", eventKey)
	}
}
