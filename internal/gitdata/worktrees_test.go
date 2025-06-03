package gitdata

import (
	"testing"
)

func TestWorktreeList(t *testing.T) {
	client := TestCLIClient()
	got, err := client.Worktrees()
	assertNoErr(t, err)
	expected := []Worktree{
		//{Name: "bare_tree", Path: "/path/to/bare_tree", Head: "", Branch: "", Bare: true}, // filtered
		{Name: "tree1", Path: "/path/to/tree1", Head: "sha-tree1", Branch: "tree1", BranchRef: "refs/heads/tree1"},
		{Name: "tree2", Path: "/path/to/tree2", Head: "sha-tree2", Branch: "tree2", BranchRef: "refs/heads/tree2"},
	}
	assertEqualWorktrees(t, got, expected)
}

func TestWorktreeMerged(t *testing.T) {
	client := TestCLIClient()
	worktreeMerged := Worktree{Name: "tree1", Head: "sha-tree1", Branch: "tree1"}
	worktreeNotMerged := Worktree{Name: "tree2", Head: "sha-tree2", Branch: "tree2"}
	assertMerged(t, client, worktreeMerged)
	assertNotMerged(t, client, worktreeNotMerged)
}
