package githubdata

import (
	"context"
	"fmt"

	"github.com/google/go-github/v71/github"
)

func (c RealClient) GetPR(branch string) (*github.PullRequest, error) {
	ctx := context.Background()

	prs, _, err := c.github.PullRequests.List(ctx, c.UpstreamOwner, c.Repo, &github.PullRequestListOptions{
		State: "all",
		Head:  fmt.Sprintf("%s:%s", c.Owner, branch),
	})
	if err != nil {
		return nil, fmt.Errorf("listing PRs: %w", err)
	}

	if len(prs) == 0 {
		return nil, nil
	}

	// pr := prs[0]
	// fmt.Printf("PR #%d: %s\n", pr.GetNumber(), pr.GetTitle())
	// fmt.Printf("URL: %s\n", pr.GetHTMLURL())
	return prs[0], nil
}

func (c RealClient) GetPRCheckStatus(pr *github.PullRequest) (*github.CombinedStatus, error) {
	ctx := context.Background()
	ref := pr.GetHead().GetSHA()
	checks, _, err := c.github.Repositories.GetCombinedStatus(ctx, c.UpstreamOwner, c.Repo, ref, nil)
	if err != nil {
		return nil, fmt.Errorf("get status checks: %w", err)
	}

	return checks, nil
}
