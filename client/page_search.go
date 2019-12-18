package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func SearchPage() (id string, content tview.Primitive) {
	var table *tview.Table
	var form *tview.Form

	// Define our grid
	grid := tview.NewGrid().
		SetRows(1).
		SetColumns(30, 0).
		SetBorders(false).
		SetGap(0, 2)

	// Define our table
	table = tview.NewTable().
		SetBorders(true)

	// Draw our table heading
	searchPageDrawTableHeading(table)

	// Define table selection logic
	table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEscape:
			app.SetFocus(form)
		case tcell.KeyEnter:
			table.SetSelectable(true, true)
		}
	}).SetSelectedFunc(func(row int, column int) {
		table.GetCell(row, column).SetTextColor(tcell.ColorRed)
		table.SetSelectable(false, false)
	})

	// Define our search form
	form = tview.NewForm().
		AddInputField("Search", "", 20, nil, nil).
		AddButton("Submit", func() {
			// Search
			//app.SetFocus(table)
		}).
		AddButton("Cancel", func() {
			pages.SwitchToPage("inbox")
		})

	// Layout for screens narrower than 100 cells (side bar are hidden).
	//grid.AddItem(chatGrid, 1, 0, 1, 2, 0, 0, true)

	// Layout for screens wider than 100 cells.
	grid.AddItem(form, 1, 0, 1, 1, 0, 100, true).
		AddItem(table, 1, 1, 1, 1, 0, 100, false)

	return "search", grid
}

func searchPageDrawTableHeading(table *tview.Table) {
	table.SetCell(0, 0, tview.NewTableCell("Username").SetTextColor(tcell.ColorPink))
}
