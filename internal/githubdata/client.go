package githubdata

import (
	"os/exec"
	"strings"

	"github.com/google/go-github/v71/github"
)

type Client interface {
	GetPR(branch string) (*github.PullRequest, error)
	GetPRCheckStatus(*github.PullRequest) (*github.CombinedStatus, error)
}

type RealClient struct {
	UpstreamOwner string
	Owner         string
	Repo          string
	github        *github.Client
}

func NewClient(upstreamOwner, owner, repo string) (*RealClient, error) {
	client := RealClient{
		UpstreamOwner: upstreamOwner,
		Owner:         owner,
		Repo:          repo,
	}
	token, err := client.ReadToken()
	if err != nil {
		return nil, err
	}
	client.github = github.NewClient(nil).WithAuthToken(token)
	return &client, nil
}

func (c *RealClient) ReadToken() (string, error) {
	cmd := exec.Command("gh", "auth", "token")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	authToken := strings.TrimRight(string(output), "\n")
	return authToken, nil
}
