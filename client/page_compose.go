package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var	(
	composeMessageContainer *tview.TextView
	composeMessageField *tview.InputField
	composeToField *tview.InputField
)

func ComposePage() (id string, content tview.Primitive) {
	grid := tview.NewGrid().
		SetRows(1).
		SetColumns(30, 0).
		SetBorders(false).
		SetGap(0, 2)

	userGrid :=  tview.NewFlex()
	chatGrid :=  tview.NewFlex()

	//userGrid.SetBorder(true)
	//userGrid.SetBorderPadding(1,1,1,1,)

	chatGrid.SetBorder(true)
	chatGrid.SetBorderPadding(1,1,1,1)
	chatGrid.SetTitle(" Compose New Message ")

	messageContainer = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	messageContainer.SetScrollable(true)
	//messageContainer.SetBorder(true)

	//SetFixed(1, 1)
	//SetDynamicColors(true).
	//SetRegions(true).
	//SetWordWrap(true).
	//SetChangedFunc(func() {
	//	app.Draw()
	//})

	inputField := tview.NewTextView().
		SetText("\n\n------\nThis is an encrypted message sent via 426c").
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			if inputField.HasFocus() {
				app.Draw()
			}
		})

	toInputField := tview.NewInputField().
		SetPlaceholder("Enter username")

	cancelButton := tview.NewButton("Cancel").SetSelectedFunc(func() {
		pages.SwitchToPage("inbox")
	})
	cancelButton.SetBorder(true).SetRect(0, 0, 0, 1)

	sendButton := tview.NewButton("Send Message").SetSelectedFunc(func() {
		pages.SwitchToPage("inbox")
	})
	sendButton.SetBorder(true).SetRect(0, 0, 0, 1)

	toInputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			app.SetFocus(inputField)
		case tcell.KeyESC:
			app.SetFocus(cancelButton)
		}
		return event
	})

	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			app.SetFocus(sendButton)
		case tcell.KeyESC:
			app.SetFocus(cancelButton)
		default:

		}
		return event
	})

	sendButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			app.SetFocus(toInputField)
		case tcell.KeyESC:
			app.SetFocus(cancelButton)
		}
		return event
	})

	cancelButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyESC:
			pages.SwitchToPage("inbox")
		}
		return event
	})

	// Layout for screens wider than 100 cells.
	grid.AddItem(userGrid, 1, 0, 1, 1, 0, 100, false).
		AddItem(chatGrid, 1, 1, 1, 1, 0, 100, true)

	userGrid.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(cancelButton, 3, 1, true), 0, 2, true)

	chatGrid.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(toInputField, 1, 1, true).
		AddItem(inputField, 0, 2, false).
		AddItem(sendButton, 3, 1, false), 0, 2, true)

	messageContainer.SetScrollable(true)

	return "compose", grid
}
