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
