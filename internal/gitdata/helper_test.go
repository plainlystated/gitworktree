package gitdata

import (
	"reflect"
	"testing"
)

func assertEqual[V comparable](t *testing.T, got, expected V) {
	t.Helper()

	if expected != got {
		t.Errorf(`assert.Equal(t, got: %v, expected: %v)`, got, expected)
	}
}

func assertEqualWorktrees(t *testing.T, got, expected []Worktree) {
	t.Helper()

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Worktree list not equal; got:\n%v\n, expected:\n%v)", got, expected)
	}
}

func assertNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func assertMerged(t *testing.T, client Service, worktree Worktree) {
	t.Helper()
	merged, err := client.IsMerged(worktree)
	assertNoErr(t, err)
	if !merged {
		t.Errorf("Expected worktree to be merged: %s", worktree.Name)
	}
}

func assertNotMerged(t *testing.T, client Service, worktree Worktree) {
	t.Helper()
	merged, err := client.IsMerged(worktree)
	assertNoErr(t, err)
	if merged {
		t.Errorf("Expected worktree to be unmerged: %s", worktree.Name)
	}
}
