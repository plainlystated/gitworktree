package main

import (
	"fmt"
	"log"

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
