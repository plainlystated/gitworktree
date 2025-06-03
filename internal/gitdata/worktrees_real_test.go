package gitdata

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

var testRepoBaseDir = "/tmp/gwttests"

func TestUnmockedClient(t *testing.T) {
	gitdir, name := setupGit()
	defer func() {
		cleanupGitRepo(testRepoBaseDir)
	}()

	client := Service{
		Client: CLIClient{
			RemoteMain: "master", // there is no remote
			Dir:        gitdir,
		},
	}

	t.Run("Worktrees", func(t *testing.T) {
		worktrees, err := client.Worktrees()
		assertNoErr(t, err)
		assertEqual(t, 3, len(worktrees))
		for i, expected := range []struct {
			name   string
			path   string
			branch string
		}{
			{name, gitdir, "master"},
			{fmt.Sprintf("%s_wt1", name), fmt.Sprintf("%s_wt1", gitdir), "branch1"},
			{fmt.Sprintf("%s_wt2", name), fmt.Sprintf("%s_wt2", gitdir), "branch2"},
		} {
			assertEqual(t, expected.name, worktrees[i].Name)
			assertEqual(t, expected.path, worktrees[i].Path)
			assertEqual(t, 40, len(worktrees[i].Head))
			assertEqual(t, expected.branch, worktrees[i].Branch)
			if !worktrees[i].UpdatedAt.After(time.Now().Add(-1 * time.Minute)) {
				t.Errorf("expected %#v to be UpdatedAt in the past minute", worktrees[i])
			}
			if !worktrees[i].UpdatedAt.Before(time.Now().Add(1 * time.Second)) {
				t.Errorf("expected %#v to be UpdatedAt before now", worktrees[i])
			}
		}
	})

	t.Run("IsMerged", func(t *testing.T) {
		worktrees, err := client.Worktrees()
		assertNoErr(t, err)

		// WT 1 is a fresh branch, so technically "merged"
		assertMerged(t, client, worktrees[1])

		// WT 2 has an extra commit beyond master
		assertNotMerged(t, client, worktrees[2])
	})
}

func runCmd(dir string, name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("command failed: %s %v\nError: %v", name, args, err)
	}
}

func setupGit() (string, string) {
	randomDir := fmt.Sprintf("repo_%d", rand.Intn(1000000))
	repoPath := filepath.Join(testRepoBaseDir, randomDir)

	if err := os.MkdirAll(repoPath, 0755); err != nil {
		log.Fatalf("failed to create repo directory: %v", err)
	}

	runCmd(repoPath, "git", "init", "--quiet")
	runCmd(repoPath, "git", "config", "user.email", "gwt-test@example.com")

	createDummyFile(repoPath, "README.md", true)

	worktree1 := filepath.Join(repoPath, "..", randomDir+"_wt1")
	worktree2 := filepath.Join(repoPath, "..", randomDir+"_wt2")
	runCmd(repoPath, "git", "worktree", "add", worktree1, "-b", "branch1", "--quiet")
	runCmd(repoPath, "git", "worktree", "add", worktree2, "-b", "branch2", "--quiet")

	createDummyFile(repoPath, "master-branch-growing", true)
	createDummyFile(worktree2, "tree2-growing", true)

	return repoPath, randomDir
}

func createDummyFile(dir, filename string, commit bool) {
	dummyFile := filepath.Join(dir, filename)
	if err := os.WriteFile(dummyFile, []byte("Some file contents...\n"), 0644); err != nil {
		log.Fatalf("failed to write file: %v", err)
	}
	if commit {
		runCmd(dir, "git", "add", filename)
		runCmd(dir, "git", "commit", "-m", filename, "--quiet")
	}
}

func cleanupGitRepo(dir string) {
	if err := os.RemoveAll(dir); err != nil {
		log.Printf("failed to delete repo: %v", err)
	}
}
