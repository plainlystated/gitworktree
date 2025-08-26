package gitdata

type Service struct {
	Client Client
}

func (s Service) Worktrees() ([]Worktree, error) {
	return s.Client.Worktrees()
}

func (s Service) IsMerged(w Worktree) (bool, error) {
	return s.Client.IsMerged(w)
}

func (s Service) DeleteWorktree(w Worktree) error {
	return s.Client.DeleteWorktree(w)
}
