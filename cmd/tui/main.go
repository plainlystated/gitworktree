package main

import (
	"fmt"
	"log"
	"os"

	"github.com/plainlystated/gitworktree/internal/gitdata"
	"github.com/plainlystated/gitworktree/internal/githubdata"
	"github.com/rivo/tview"
)

type app struct {
	tview  *tview.Application
	git    gitdata.Service
	github githubdata.Service
}

func main() {
	if len(os.Args) == 1 {
		tui()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "add":
		if len(os.Args) != 3 {
			fmt.Println("Expected 1 arg: branch_name")
			os.Exit(1)
		}
		CreateWorktree(os.Args[2])

	case "co":
		if len(os.Args) != 3 {
			fmt.Println("Expected 1 arg: branch_name")
			os.Exit(1)
		}
		CheckoutWorktree(os.Args[2])
	case "list":
		if len(os.Args) != 2 {
			fmt.Println("Expected no args")
			os.Exit(1)
		}
		ListWorktrees()

	default:
		fmt.Println("Expected one of: add, co, list")
		os.Exit(1)
	}
}

func tui() {
	githubService, err := githubdata.NewService("lewiscountypress", "plainlystated", "workbench")
	if err != nil {
		log.Fatal(err.Error())
	}
	app := app{
		tview: tview.NewApplication(),
		// gh:    gitdata.DefaultCLIClient(),
		// gh: gitdata.TestCLIClient(),
		git: gitdata.Service{
			Client: gitdata.CLIClient{
				RemoteMain: "upstream/master",
				// Exec:       gitdata.LocalCLIExec{},
			},
		},
		github: githubService,
	}

	worktreeList := app.worktreeList()

	// if err := app.SetRoot(box, true).Run(); err != nil {
	if err := app.tview.SetRoot(worktreeList, true).Run(); err != nil {
		panic(err)
	}
}

func (a app) errorView(err error) tview.Primitive {
	textview := tview.NewTextView().
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			a.tview.Draw()
		})
	textview.SetBorder(true).
		SetTitle("Worktrees")
	fmt.Fprintf(textview, "%s", err.Error())

	return textview
}
