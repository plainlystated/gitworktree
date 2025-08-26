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

func WorktreeByName(wts []Worktree, name string) (wt Worktree, found bool) {
	for _, wt := range wts {
		if wt.Name == name {
			return wt, true
		}
	}
	return Worktree{}, false
}
