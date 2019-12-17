package main

import "github.com/rivo/tview"

var	(
	messageContainer *tview.TextView
	userListContainer *tview.List
)

func InboxPage() (id string, content tview.Primitive) {
	grid := tview.NewGrid().
		SetRows(1).
		SetColumns(30, 0).
		SetBorders(true).
		SetGap(0, 2)


	return "inbox", grid
}
