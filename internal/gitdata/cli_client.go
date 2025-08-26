package gitdata

// There is a "real" go git library, but:
// - it has poor support for worktrees
// - it has a TON of dependencies; seems like overkill for our limited needs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type CLIClient struct {
	Dir        string
	RemoteMain string
}

func DefaultService() Service {
	return Service{
		Client: CLIClient{
			RemoteMain: "origin/master",
		},
	}
}

func (c CLIClient) IsMerged(worktree Worktree) (bool, error) {
	mergeBase, err := c.mergeBase(c.RemoteMain, worktree.Branch)
	if err != nil {
		return false, err
	}

	if mergeBase == worktree.Head {
		return true, nil
	}
	return false, nil
}

func (c CLIClient) DeleteWorktree(wt Worktree) error {
	err := os.RemoveAll(wt.Path)
	if err != nil {
		return fmt.Errorf("deleting %s: %w", wt.Path, err)
	}

	cmd := exec.Command("git", "worktree", "prune")
	cmd.Dir = filepath.Join(wt.Path, "..")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git prune: %w, %s", err, output)
	}
	return nil
}

func (c CLIClient) Worktrees() ([]Worktree, error) {
	bytes, err := c.worktreeList()
	str := string(bytes)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(str, "\n")
	trees := []Worktree{}
	wip := Worktree{}
	for _, line := range lines {
		if line == "" {
			if !wip.Bare && wip != (Worktree{}) {
				trees = append(trees, wip)
			}
			wip = Worktree{}
		}
		parts := strings.Split(line, " ")
		switch parts[0] {
		case "bare":
			wip.Bare = true
		case "worktree":
			wip.Path = parts[1]
			pathParts := strings.Split(parts[1], "/")
			wip.Name = pathParts[len(pathParts)-1]
		case "HEAD":
			wip.Head = CommitSHA(parts[1])

			wip.UpdatedAt, err = c.commitTime(wip.Head)
			if err != nil {
				return nil, fmt.Errorf("getting timestamp for commit %s: %w", wip.Head, err)
			}
		case "branch":
			wip.BranchRef = parts[1]
			pathParts := strings.Split(parts[1], "/")
			wip.Branch = pathParts[len(pathParts)-1]
		}
	}

	return trees, nil
}

func (c CLIClient) commitTime(commit CommitSHA) (time.Time, error) {
	cmd := exec.Command("git", "log", "-n", "1", "--format=%cd", "--date=iso-strict", string(commit))
	cmd.Dir = c.Dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %s", err, output)
	}
	t, err := time.Parse(time.RFC3339, strings.TrimSpace(string(output)))
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %s", err, output)
	}
	return t, nil
}

func (c CLIClient) worktreeList() ([]byte, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	cmd.Dir = c.Dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, fmt.Errorf("%w: %s", err, output)
	}
	return output, nil
}

func (c CLIClient) mergeBase(commit1, commit2 string) (CommitSHA, error) {
	cmd := exec.Command("git", "merge-base", commit1, commit2)
	cmd.Dir = c.Dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return CommitSHA(string(output)), fmt.Errorf("%w: %s", err, output)
	}
	return CommitSHA(strings.TrimRight(string(output), "\n")), nil
}
