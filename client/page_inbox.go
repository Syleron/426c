package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/olekukonko/tablewriter"
	"github.com/rivo/tview"
	"github.com/syleron/426c/common/models"
	"github.com/syleron/426c/common/packet"
	"github.com/syleron/426c/common/utils"
)

var (
	userListContainer *tview.Table
	inboxToField *tview.InputField
)

func InboxPage() (id string, content tview.Primitive) {
	var inputField *tview.InputField
	var selectedUsername string

	grid := tview.NewGrid().
		SetRows(1).
		SetColumns(20, 0).
		SetBorders(false).
		SetGap(0, 1)

	userGrid := tview.NewFlex()
	chatGrid := tview.NewFlex()

	chatGrid.SetBorder(false)

	userGrid.SetBorderPadding(0, 1, 0, 0)
	chatGrid.SetBorderPadding(0, 1, 0, 0)

	messageContainer := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	messageContainer.SetScrollable(true)

	userListContainer = tview.NewTable()
	userListContainer.SetBorder(false)
	userListContainer.SetBorderPadding(1, 0, 0, 0)
	userListContainer.SetSelectable(true, true)

	userListContainer.
		SetSelectedFunc(func(row, column int) {
			username := userListContainer.GetCell(row, column)
			// Mark our selected left table cell
			username.SetTextColor(tcell.ColorRed)
			// Set our selected username
			selectedUsername = username.Text
			// Load our messages for the user
			loadMessages(selectedUsername, messageContainer)
			// Set focus on our message container
			app.SetFocus(inputField)
		},
	)

	inputField = tview.NewInputField().
		SetPlaceholder("Send message...").
		//SetAcceptanceFunc(tview.InputFieldInteger).
		SetDoneFunc(func(key tcell.Key) {
			switch key {
			case tcell.KeyESC:
				app.SetFocus(userListContainer)
				// reset selection
				for i := 0; i < userListContainer.GetRowCount(); i++ {
					userListContainer.GetCell(i, 0).SetTextColor(tcell.ColorWhite)
				}
			case tcell.KeyUp:
				r, c := messageContainer.GetScrollOffset()
				messageContainer.ScrollTo(r - 1, c)
			case tcell.KeyDown:
				r, c := messageContainer.GetScrollOffset()
				messageContainer.ScrollTo(r + 1, c)
			case tcell.KeyEnter:
				if inputField.GetText() == "" {
					return
				}
				// submit our message
				submitMessage(selectedUsername, inputField.GetText())
				// clear out our input
				inputField.SetText("")
				// reload our messages
				loadMessages(selectedUsername, messageContainer)
			}
		})

	composeButton := tview.NewButton("Compose").SetSelectedFunc(func() {
		creditBlocks(12)
		pages.SwitchToPage("compose")
	})
	composeButton.SetBorder(true).SetRect(0, 0, 0, 1)

	inboxToField = tview.NewInputField().
		SetPlaceholder("Search user")

	composeButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyESC:
			inboxQuitModal()
		case tcell.KeyTAB:
			app.SetFocus(userListContainer)
		}
		return event
	})

	userListContainer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			app.SetFocus(inboxToField)
		}
		return event
	})

	inboxToField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			app.SetFocus(userListContainer)
		case tcell.KeyESC:
			inboxQuitModal()
		case tcell.KeyEnter:
			// Check if user exists and get public key details
			_, err := client.Send(packet.CMD_USER, utils.MarshalResponse(&models.UserRequestModel{
				Username: inboxToField.GetText(),
			}))
			if err != nil {
				panic(err)
			}
		}
		return event
	})

	// Layout for screens wider than 100 cells.
	grid.AddItem(userGrid, 1, 0, 1, 1, 0, 50, true).
		AddItem(chatGrid, 1, 1, 1, 1, 0, 100, false)

	userGrid.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(inboxToField, 1, 1, true).
		AddItem(userListContainer, 0, 1, true), 0, 2, true)

	chatGrid.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(messageContainer, 0, 1, false).
		AddItem(inputField, 1, 1, false), 0, 2, false)

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
	for i, user := range users {
		if user.Username != lUser {
			userListContainer.SetCell(i, 0, tview.NewTableCell(user.Username))
		}
	}
}

func loadMessages(username string, container *tview.TextView) {
	// clear our messages
	container.Clear()
	// Get our messages
	messages, _ := dbMessagesGet(username, lUser)
	reverseAny(messages)
	for _, message := range messages {
		var fmsg string
		var color string
		color = "[gray]"
		// Set our message stats
		if !message.Success {
			fmsg += "[red]! " + color
		} else {
			fmsg += color
		}
		// Set our time
		fmsg += message.Date.Format("15:04:05")
		// Set from/to
		if message.To == lUser {
			fmsg += " <[white]" + message.To + color +  "> "
		} else {
			fmsg += " <[darkcyan]" + message.From + color + "> "
		}
		// Set our message
		if message.To == lUser {
			fmsg += decryptMessage(message.ToMessage)
		} else {
			fmsg += decryptMessage(message.FromMessage)
		}

		fmt.Fprintf(container,`%s %v`, fmsg, tablewriter.NEWLINE)
	}
	container.ScrollToEnd()
}

func inboxQuitModal() {
	showModal(ClientModal{
		Message:      "Are you sure you would like to disconnect from the 426c network?",
		SubmitButton: "Exit",
		CancelButton: "Cancel",
		Continue: func() {
			app.Stop()
		},
	})
}