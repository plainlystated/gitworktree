package main

import (
	"fmt"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/plainlystated/gitworktree/internal/gitdata"
	"github.com/rivo/tview"
)

func (a app) buildWorktreeList() *tview.Table {
	table := tview.NewTable().
		SetFixed(1, 2).
		SetSelectable(true, false)
	return table
}

func (a app) refreshWorktreeList() error {
	a.worktreeList.Clear()

	worktrees, err := a.git.Worktrees()
	if err != nil {
		return fmt.Errorf("error fetching worktrees: %w", err)
	}

	sort.Slice(worktrees, func(i, j int) bool {
		return worktrees[i].UpdatedAt.After(worktrees[j].UpdatedAt)
	})

	name := tview.NewTableCell("Name").
		SetTextColor(tcell.ColorYellow).
		SetSelectable(false)
	path := tview.NewTableCell("Merged").
		SetTextColor(tcell.ColorYellow).
		SetSelectable(false)
	pr := tview.NewTableCell("PR").
		SetTextColor(tcell.ColorYellow).
		SetSelectable(false)
	a.worktreeList.SetCell(0, 0, name)
	a.worktreeList.SetCell(0, 1, path)
	a.worktreeList.SetCell(0, 2, pr)

	for i, worktree := range worktrees {
		name := a.nameCell(worktree)
		merged := a.mergedCell(worktree)

		a.worktreeList.SetCell(i+1, 0, name)
		a.worktreeList.SetCell(i+1, 1, merged)

		go func() {
			pr, isMerged := a.prCell(worktree)
			a.tview.QueueUpdateDraw(func() {
				if isMerged {
					name.SetTextColor(tcell.ColorDimGray)
				}
				a.worktreeList.SetCell(i+1, 2, pr)
			})
		}()

	}

	a.handleDelete(a.worktreeList)

	return nil
}

func (a app) nameCell(worktree gitdata.Worktree) *tview.TableCell {
	cell := tview.NewTableCell(worktree.Name).
		SetTextColor(tcell.ColorWhite)
	return cell
}

func (a app) mergedCell(worktree gitdata.Worktree) *tview.TableCell {
	var cell *tview.TableCell
	merged, err := a.git.IsMerged(worktree)
	if err != nil {
		cell = tview.NewTableCell(err.Error()).
			SetTextColor(tcell.ColorRed)
	} else {
		if merged {
			cell = tview.NewTableCell("✓").
				SetTextColor(tcell.ColorGreen)
		} else {
			cell = tview.NewTableCell("")
		}
	}
	cell.SetSelectable(true)
	return cell
}

func (a app) prCell(worktree gitdata.Worktree) (*tview.TableCell, bool) {
	pr, err := a.github.GetPR(worktree.Branch)
	if err != nil {
		return tview.NewTableCell(err.Error()), false
	}
	if pr == nil {
		return tview.NewTableCell(""), false
	}

	prChecksState, err := a.github.PRChecksPassed(pr)
	if err != nil {
		return tview.NewTableCell(err.Error()), false
	}

	prCheckStr := prChecksState
	color := tcell.ColorWhite
	switch prChecksState {
	case "success":
		prCheckStr = "✓"
		color = tcell.ColorGreen
	case "failure":
		prCheckStr = "❌"
		color = tcell.ColorRed
	case "pending":
		prCheckStr = "↻"
		color = tcell.ColorGreenYellow
	}

	if *pr.State == "closed" {
		color = tcell.ColorDimGray
	}

	prStr := fmt.Sprintf("#%d %s", pr.GetNumber(), prCheckStr)
	return tview.NewTableCell(prStr).
		SetTextColor(color), prChecksState == "success"
}

func (a app) handleDelete(table *tview.Table) {
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		a.clearStatus()

		if event.Rune() == 'd' {
			row, _ := table.GetSelection()
			cell := table.GetCell(row, 0)
			message := fmt.Sprintf("Delete worktree %s?", cell.Text)

			modal := tview.NewModal().
				SetText(message).
				AddButtons([]string{"OK", "Cancel"}).
				SetDoneFunc(func(buttonIndex int, buttonLabel string) {
					if buttonLabel == "OK" {
						worktrees, err := a.git.Worktrees()
						if err != nil {
							a.showStatus(fmt.Sprintf("Error: %s", err.Error()))
						}
						wt, found := gitdata.WorktreeByName(worktrees, cell.Text)
						if !found {
							a.showStatus(fmt.Sprintf("Invalid selection: %s", cell.Text))
						}
						err = a.git.DeleteWorktree(wt)
						if err != nil {
							a.showStatus(fmt.Sprintf("Error deleting worktree: %s", err.Error()))
						} else {
							a.showStatus(fmt.Sprintf("Deleted worktree and local copy: %s", cell.Text))
							a.refreshWorktreeList()
						}
					}
					a.showGrid()
				})

			a.tview.SetRoot(modal, true).SetFocus(modal)
			return nil
		}
		return event
	})
}
