package main

import (
	"fmt"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/plainlystated/gitworktree/internal/gitdata"
	"github.com/rivo/tview"
)

func (a app) worktreeList() tview.Primitive {
	worktrees, err := a.git.Worktrees()
	if err != nil {
		return a.errorView(fmt.Errorf("error fetching worktrees: %w", err))
	}

	sort.Slice(worktrees, func(i, j int) bool {
		return worktrees[i].UpdatedAt.After(worktrees[j].UpdatedAt)
	})

	table := tview.NewTable().
		SetFixed(len(worktrees), 2).
		SetSelectable(true, false)
	name := tview.NewTableCell("Name").
		SetTextColor(tcell.ColorYellow).
		SetSelectable(false)
	path := tview.NewTableCell("Merged").
		SetTextColor(tcell.ColorYellow).
		SetSelectable(false)
	pr := tview.NewTableCell("PR").
		SetTextColor(tcell.ColorYellow).
		SetSelectable(false)
	table.SetCell(0, 0, name)
	table.SetCell(0, 1, path)
	table.SetCell(0, 2, pr)

	for i, worktree := range worktrees {
		name := a.nameCell(worktree)
		merged := a.mergedCell(worktree)

		table.SetCell(i+1, 0, name)
		table.SetCell(i+1, 1, merged)

		go func() {
			pr, isMerged := a.prCell(worktree)
			a.tview.QueueUpdateDraw(func() {
				if isMerged {
					name.SetTextColor(tcell.ColorDimGray)
				}
				table.SetCell(i+1, 2, pr)
			})
		}()

	}

	return table
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
