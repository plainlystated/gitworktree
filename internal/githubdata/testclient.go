package githubdata

import (
	"github.com/google/go-github/v71/github"
)

type TestClient struct{}

func (c TestClient) GetPR(branch string) (*github.PullRequest, error) {
	return &github.PullRequest{}, nil
}

func (c TestClient) GetPRCheckStatus(pr *github.PullRequest) (*github.CombinedStatus, error) {
	return nil, nil
}
