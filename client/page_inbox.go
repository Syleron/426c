package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/syleron/426c/common/models"
	"github.com/syleron/426c/common/packet"
	"github.com/syleron/426c/common/utils"
	"sync"
)

var (
	inboxMessageContainer *tview.TextView
	userListContainer *tview.Table
	inboxToField *tview.InputField
	inboxMessageContainerLock sync.Mutex
	inboxSelectedUsername string
	inboxFailedMessageCount int
)

// TODO: Upon receiving a message from a user that you dont already have, nothing shows.
// TODO: Take into account multiple accounts sending messages to the same user


func InboxPage() (id string, content tview.Primitive) {
	var inputField *tview.InputField

	grid := tview.NewGrid().
		SetRows(1).
		SetColumns(15, 0).
		SetBorders(false).
		SetGap(0, 1)

	userGrid := tview.NewFlex()
	chatGrid := tview.NewFlex()

	chatGrid.SetBorder(false)

	userGrid.SetBorderPadding(0, 1, 0, 0)
	chatGrid.SetBorderPadding(0, 1, 0, 0)

	inboxMessageContainer = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	inboxMessageContainer.SetScrollable(true)

	userListContainer = tview.NewTable()
	userListContainer.SetBorder(false)
	userListContainer.SetBorderPadding(1, 0, 0, 0)
	userListContainer.SetSelectable(true, true)

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
				r, c := inboxMessageContainer.GetScrollOffset()
				inboxMessageContainer.ScrollTo(r - 1, c)
			case tcell.KeyDown:
				r, c := inboxMessageContainer.GetScrollOffset()
				inboxMessageContainer.ScrollTo(r + 1, c)
			case tcell.KeyEnter:
				if inputField.GetText() == "" {
					return
				}
				// submit our message
				messageSubmit(inboxSelectedUsername, inputField.GetText())
				// clear out our input
				inputField.SetText("")
			}
		})

	inboxToField = tview.NewInputField().
		SetPlaceholder("Search user")

	userListContainer.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyESC:
			inboxQuitModal()
		case tcell.KeyTAB:
			app.SetFocus(inboxToField)
		case tcell.KeyEnter:
			row, column := userListContainer.GetOffset()
			username := userListContainer.GetCell(row, column)
			// Mark our selected left table cell
			username.SetTextColor(tcell.ColorRed)
			// Set our selected username
			inboxSelectedUsername = username.Text
			// Load our messages for the user
			go messageLoad(inboxSelectedUsername, inboxMessageContainer)
			go inboxRetryFailedMessages(inboxSelectedUsername)
			// Set focus on our message container
			app.SetFocus(inputField)
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
		AddItem(inboxMessageContainer, 0, 1, false).
		AddItem(inputField, 1, 1, false), 0, 2, false)

	// Get our contacts
	go drawContactsList()

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
	count := 0 // custom counter as skipping our user screws with the for one
	for _, user := range users {
		if user.Username != client.Username {
			userListContainer.SetCell(count, 0, tview.NewTableCell(user.Username))
			count++
		}
	}
	app.Draw()
}


func inboxRetryFailedMessages(username string) {
	// reset counter
	inboxFailedMessageCount = 0
	// Get our messages
	messages, _ := dbMessagesGet(username, client.Username)
	reverseAny(messages)
	for _, message := range messages {
		if !message.Success {
			client.cmdMsgTo(&message) // retry message
		}
	}
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
