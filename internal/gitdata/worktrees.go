package gitdata

import "time"

type Worktree struct {
	Name      string
	Path      string
	Head      CommitSHA
	Branch    string
	BranchRef string
	Bare      bool
	UpdatedAt time.Time
}
