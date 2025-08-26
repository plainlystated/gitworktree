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
	tview        *tview.Application
	git          gitdata.Service
	github       githubdata.Service
	grid         *tview.Grid
	worktreeList *tview.Table
	footer       *tview.TextView
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

	footer := tview.NewTextView().
		SetTextAlign(tview.AlignCenter)

	grid := tview.NewGrid().SetRows(3, 0, 3).
		AddItem(tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText("gitworktree"),
			0, 0, 1, 1, 0, 0, false).
		AddItem(footer,
			2, 0, 1, 1, 0, 0, false)

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
		grid:   grid,
		footer: footer,
		github: githubService,
	}
	app.worktreeList = app.buildWorktreeList()
	app.refreshWorktreeList()

	grid.AddItem(app.worktreeList, 1, 0, 1, 1, 0, 0, true)

	if err := app.showGrid().Run(); err != nil {
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
