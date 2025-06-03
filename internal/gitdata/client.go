package gitdata

type Client interface {
	Worktrees() ([]Worktree, error)
	IsMerged(Worktree) (bool, error)
}
