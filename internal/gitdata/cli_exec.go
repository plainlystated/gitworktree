package gitdata

type CommitSHA string

// type cliExec interface {
// 	WorktreeList() ([]byte, error)
// 	MergeBase(commitish1, commitish2 string) (CommitSHA, error)
// 	// CommitSHA(branch string) (CommitSHA, error)
// }
//
// type LocalCLIExec struct {
// 	Dir string
// }

// func (c LocalCLIExec) WorktreeList() ([]byte, error) {
// 	cmd := exec.Command("git", "worktree", "list", "--porcelain")
// 	cmd.Dir = c.Dir
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return output, fmt.Errorf("%w: %s", err, output)
// 	}
// 	return output, nil
// }

// func (c LocalCLIExec) MergeBase(commit1, commit2 string) (CommitSHA, error) {
// 	cmd := exec.Command("git", "merge-base", commit1, commit2)
// 	cmd.Dir = c.Dir
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return CommitSHA(string(output)), fmt.Errorf("%w: %s", err, output)
// 	}
// 	return CommitSHA(strings.TrimRight(string(output), "\n")), nil
// }

// func (c LocalCLIExec) CommitSHA(commitish string) (CommitSHA, error) {
// 	cmd := exec.Command("git", "rev-parse", commitish)
// 	cmd.Dir = c.Dir
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return CommitSHA(string(output)), fmt.Errorf("%w: %s", err, output)
// 	}
// 	return CommitSHA(strings.TrimRight(string(output), "\n")), nil
// }
