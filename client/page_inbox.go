package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var (
	messageContainer  *tview.TextView
	userListContainer *tview.Table
	inputField        *tview.InputField
)

func InboxPage() (id string, content tview.Primitive) {
	grid := tview.NewGrid().
		SetRows(1).
		SetColumns(30, 0).
		SetBorders(false).
		SetGap(0, 2)

	userGrid := tview.NewFlex()
	chatGrid := tview.NewFlex()

	chatGrid.SetBorder(true)
	chatGrid.SetBorderPadding(1, 1, 1, 1)

	messageContainer = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	messageContainer.SetScrollable(true)
	//messageContainer.SetBorder(true)

	userListContainer = tview.NewTable()
	userListContainer.SetBorder(true)
	userListContainer.SetBorderPadding(1, 1, 1, 1)

	inputField = tview.NewInputField().
		SetPlaceholder("Send message...").
		//SetAcceptanceFunc(tview.InputFieldInteger).
		SetDoneFunc(func(key tcell.Key) {
			switch key {
			case tcell.KeyUp:
				//r, c := chatContainer.GetScrollOffset()
				//chatContainer.ScrollTo(r - 1, c)
			case tcell.KeyDown:
				//r, c := chatContainer.GetScrollOffset()
				//chatContainer.ScrollTo(r + 1, c)
			case tcell.KeyEnter:
				//inputField.SetText("")
			}
		})

	composeButton := tview.NewButton("Compose").SetSelectedFunc(func() {
		creditBlocks(12)
		pages.SwitchToPage("compose")
	})
	composeButton.SetBorder(true).SetRect(0, 0, 0, 1)

	composeButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyESC:
			showModal(ClientModal{
				Message:      "Are you sure you would like to disconnect from the 426c network?",
				SubmitButton: "Exit",
				CancelButton: "Cancel",
				Continue: func() {
					app.Stop()
				},
			})
		case tcell.KeyTAB:
			app.SetFocus(userListContainer)
		}
		return event
	})

	userListContainer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			app.SetFocus(composeButton)
		}
		return event
	})

	// Layout for screens wider than 100 cells.
	grid.AddItem(userGrid, 1, 0, 1, 1, 0, 100, true).
		AddItem(chatGrid, 1, 1, 1, 1, 0, 100, false)

	userGrid.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(composeButton, 3, 1, true).
		AddItem(userListContainer, 0, 1, true), 0, 2, true)

	chatGrid.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(messageContainer, 0, 1, false).
		AddItem(inputField, 1, 1, false), 0, 2, false)

	messageContainer.SetScrollable(true)

	// Get our contacts
	drawContactsList()

	return "inbox", grid
}

func drawContactsList() {
	// Clear our current list
	userListContainer.Clear()
	// Get our contacts from our db
	users, err := dbUserList()
	if err != nil {
		app.Stop()
	}
	// List all of our contacts in our local DB
	for _, user := range users {
		userListContainer.SetCell(0, 0, tview.NewTableCell(user.Username))
	}
}
