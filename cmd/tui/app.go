package main

import "github.com/rivo/tview"

func (a app) showStatus(msg string) {
	a.footer.SetText(msg)
}

func (a app) clearStatus() {
	a.footer.SetText("")
}

func (a app) showGrid() *tview.Application {
	return a.
		tview.
		SetRoot(a.grid, true).
		SetFocus(a.grid)
}
