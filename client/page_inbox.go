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
		SetBorders(false).
		SetGap(0, 2)

	userGrid :=  tview.NewFlex()
	chatGrid :=  tview.NewFlex()

	//userGrid.SetBorder(true)
	//userGrid.SetBorderPadding(1,1,1,1,)

	chatGrid.SetBorder(true)
	chatGrid.SetBorderPadding(1,1,1,1)

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
	userListContainer.SetBorderPadding(1,1,1,1)

		//SetFixed(1, 1)
		//SetDynamicColors(true).
		//SetRegions(true).
		//SetWordWrap(true).
		//SetChangedFunc(func() {
		//	app.Draw()
		//})

	userListContainer.SetCell(0, 0, tview.NewTableCell("Willifer (Online)"))
	userListContainer.SetCell(1, 0, tview.NewTableCell("Haroto (Offline)"))

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
		creditBlocks(12)
		pages.SwitchToPage("search")
	})
	button.SetBorder(true).SetRect(0, 0, 0, 1)

	// Layout for screens wider than 100 cells.
	grid.AddItem(userGrid, 1, 0, 1, 1, 0, 100, true).
		AddItem(chatGrid, 1, 1, 1, 1, 0, 100, false)

	userGrid.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(button, 3, 1, true).
		AddItem(userListContainer, 0, 1, true), 0, 2, true)

	chatGrid.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(messageContainer, 0, 1, false).
		AddItem(inputField, 1, 1, false), 0, 2, false)

	messageContainer.SetScrollable(true)

	return "inbox", grid
}
