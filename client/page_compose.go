package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var	(
	composeMessageContainer *tview.TextView
	composeUserListContainer *tview.Table
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

	composeMessageContainer = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	composeMessageContainer.SetScrollable(true)
	//messageContainer.SetBorder(true)

	composeMessageField = tview.NewInputField().
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

	button := tview.NewButton("Cancel").SetSelectedFunc(func() {
		pages.SwitchToPage("inbox")
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

	composeMessageContainer.SetScrollable(true)

	return "compose", grid
}
