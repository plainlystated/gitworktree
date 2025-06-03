package main

import (
	"fmt"
	"log"

	"github.com/plainlystated/gitworktree/internal/githubdata"
	"github.com/rivo/tview"
)

func main() {
	githubService, err := githubdata.NewService("lewiscountypress", "plainlystated", "workbench")
	if err != nil {
		log.Fatal(err.Error())
	}

	pr, err := githubService.GetPR("sitestatus2")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(githubService.PRChecksPassed(pr))
	// client := github.NewClient(nil)

	// list all organizations for user "willnorris"
	// orgs, _, err := client.Organizations.List(context.Background(), "willnorris", nil)
	// spew.Dump(orgs)
}

func tui() {
	app := tview.NewApplication()
	textview := tview.NewTextView().
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	textview.SetBorder(true).
		SetTitle("Worktrees")
	fmt.Fprintln(textview, "asdfasdf")
	// box := tview.NewBox().SetBorder(true).SetTitle("Hello, world!")
	if err := app.SetRoot(textview, true).Run(); err != nil {
		// if err := app.SetRoot(box, true).Run(); err != nil {
		panic(err)
	}
}
