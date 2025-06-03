package githubdata

import "testing"

func TestPullRequests(t *testing.T) {
	// service := setupTestService()
	// got, err := service.GetPR("branch1")
	// assertNoErr(t, err)
	// assertEqual(t, got, 4)
}

func setupTestService() Service {
	return Service{client: TestClient{}}
}

func NewTestService(owner string) Service {
	return Service{
		client: TestClient{},
	}
}
