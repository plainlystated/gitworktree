package gitdata

type testClient struct{}

// func (c testClient) WorktreeList() ([]byte, error) {
// 	return []byte(`worktree /path/to/bare_tree
// bare
//
// worktree /path/to/tree1
// HEAD sha-tree1
// branch refs/heads/tree1
//
// worktree /path/to/tree2
// HEAD sha-tree2
// branch refs/heads/tree2
// `), nil
// }

func (c testClient) IsMerged(worktree Worktree) (bool, error) {
	switch worktree.Branch {
	case "tree1":
		return true, nil
	}
	return false, nil
}

//	func (c testClient) MergeBase(commit1, commit2 string) (CommitSHA, error) {
//		// commit1 is master, but ignored here
//		var sha CommitSHA
//		switch commit2 {
//		case "tree1":
//			sha = "sha-tree1" // tree1 is empty, ie "merged"
//		case "tree2":
//			sha = "some-prev-SHA" // tree2 has commits, so its common ancestor with master is in the past...
//		}
//		return sha, nil
//	}
func (c testClient) Worktrees() ([]Worktree, error) {
	return []Worktree{
		//{Name: "bare_tree", Path: "/path/to/bare_tree", Head: "", Branch: "", Bare: true}, // filtered
		{Name: "tree1", Path: "/path/to/tree1", Head: "sha-tree1", Branch: "tree1", BranchRef: "refs/heads/tree1"},
		{Name: "tree2", Path: "/path/to/tree2", Head: "sha-tree2", Branch: "tree2", BranchRef: "refs/heads/tree2"},
	}, nil
}

func (c testClient) DeleteWorktree(wt Worktree) error {
	return nil
}

// func (c testCLIExec) CommitSHA(commitish string) (CommitSHA, error) {
// 	var sha CommitSHA
// 	switch commitish {
// 	case "origin/master":
// 		sha = "sha-tree1" // master HEAD
// 	case "tree1":
// 		sha = "sha-tree1" // empty branch, same as master
// 	case "tree2":
// 		sha = "sha-tree2" // this branch has commits, so there's some new SHA
// 	}
// 	return sha, nil
// }

func TestCLIClient() Service {
	return Service{
		Client: testClient{},
	}
}
