package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var	(
	messageContainer *tview.TextView
	userListContainer *tview.Table
	inputField *tview.InputField
)

func InboxPage() (id string, content tview.Primitive) {
	grid := tview.NewGrid().
		SetRows(1).
		SetColumns(30, 0).
		SetBorders(true).
		SetGap(0, 2)

	userGrid :=  tview.NewFlex()
	chatGrid :=  tview.NewFlex()

	messageContainer = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	messageContainer.SetScrollable(true)

	userListContainer = tview.NewTable().
		SetFixed(1, 1)
		//SetDynamicColors(true).
		//SetRegions(true).
		//SetWordWrap(true).
		//SetChangedFunc(func() {
		//	app.Draw()
		//})

	userListContainer.SetCell(0, 0, tview.NewTableCell("testing"))
	userListContainer.SetCell(1, 0, tview.NewTableCell("testing"))
	userListContainer.SetCell(2, 0, tview.NewTableCell("testing"))

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
				//sockets.Emit(&common.Message{
				//	EventName: "chat",
				//	Data:      []byte(`{"message":"` + inputField.GetText() + `", "channel": "general"}`),
				//})
				//inputField.SetText("")
			}
		})

	button := tview.NewButton("Compose").SetSelectedFunc(func() {
	})
	button.SetBorder(true).SetRect(0, 0, 22, 3)

	// Layout for screens narrower than 100 cells (side bar are hidden).
	//grid.AddItem(chatGrid, 1, 0, 1, 2, 0, 0, true)

	// Layout for screens wider than 100 cells.
	grid.AddItem(userGrid, 1, 0, 1, 1, 0, 100, true).
		AddItem(chatGrid, 1, 1, 1, 1, 0, 100, false)

	userGrid.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(button, 0, 1, false).
		AddItem(userListContainer, 1, 1, true), 0, 2, true)

	chatGrid.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(messageContainer, 0, 1, false).
		AddItem(inputField, 1, 1, true), 0, 2, true)

	messageContainer.SetScrollable(true)

	return "inbox", grid
}
