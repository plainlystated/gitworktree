package githubdata

import "github.com/google/go-github/v71/github"

type Service struct {
	client Client
}

func NewService(upstreamOwner, owner, repo string) (Service, error) {
	client, err := NewClient(
		upstreamOwner,
		owner,
		repo,
	)
	if err != nil {
		return Service{}, err
	}
	service := Service{
		client: client,
	}
	return service, nil
}

func (s Service) GetPR(branch string) (*github.PullRequest, error) {
	return s.client.GetPR(branch)
}

func (s Service) PRChecksPassed(pr *github.PullRequest) (string, error) {
	checks, err := s.client.GetPRCheckStatus(pr)
	if err != nil {
		return "", err
	}
	return checks.GetState(), nil
}
