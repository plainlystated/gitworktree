package gitdata

import (
	"fmt"
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

func assertIncludesWorktreeName(t *testing.T, worktrees []Worktree, exp string) {
	t.Helper()
	for _, wt := range worktrees {
		if wt.Name == exp {
			return
		}
		fmt.Println(wt.Name)
	}
	t.Errorf("Expected worktrees to include %s", exp)
}

func assertExcludesWorktreeName(t *testing.T, worktrees []Worktree, exp string) {
	t.Helper()
	for _, wt := range worktrees {
		if wt.Name == exp {
			t.Errorf("Expected worktrees to not include %s", exp)
			return
		}
	}
}
