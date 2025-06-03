package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

func ListWorktrees() {
	cliExec("git", "worktree", "list")
}

func CheckoutWorktree(name string) {
	mainDir, err := mainGitDir()
	if err != nil {
		fmt.Printf("Error getting main git dir in path: %s", err)
	}
	cliExec("git", "fetch", "upstream")
	cliExec("git", "worktree", "add", fmt.Sprintf("../%s", name), name)
	cliExec("tms", fmt.Sprintf("%s/../%s", mainDir, name))
}

func CreateWorktree(name string) {
	mainDir, err := mainGitDir()
	if err != nil {
		fmt.Printf("Error getting main git dir in path: %s", err)
	}
	cliExec("git", "fetch", "upstream")
	cliExec("git", "worktree", "add", "-b", name, fmt.Sprintf("../%s", name))
	cliExec("tms", fmt.Sprintf("%s/../%s", mainDir, name))
}

func mainGitDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`^.*\.git(?:/|$)`)
	match := re.FindString(cwd)
	if match == "" {
		return "", fmt.Errorf("no repo.git dir found in current path")
	} else {
		fmt.Println("Found .git path:", match)
		return filepath.Join(match, "main"), nil
	}
}

func cliExec(bin string, args ...string) {
	cmd := exec.Command(bin, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n", output)
		fmt.Printf("Err: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(output))
}
