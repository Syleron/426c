package main

import (
	gopenpgp "github.com/ProtonMail/gopenpgp/crypto"
	"github.com/gdamore/tcell"
	"github.com/olekukonko/tablewriter"
	"github.com/rivo/tview"
	"github.com/syleron/426c/common/models"
	"github.com/syleron/426c/common/packet"
	"github.com/syleron/426c/common/utils"
	"strings"
	"sync"
	"time"
)

var (
	inboxMessageContainer *tview.TextView
	userListContainer *tview.Table
	inboxToField *tview.InputField
	inboxMessageContainerLock sync.Mutex
	inboxSelectedUsername string
	inboxFailedMessageCount int
)

// TODO: When searching and finding a user, the userlist does not reload
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
				submitMessage(inboxSelectedUsername, inputField.GetText())
				// clear out our input
				inputField.SetText("")
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
		case tcell.KeyEnter:
			row, column := userListContainer.GetOffset()
			username := userListContainer.GetCell(row, column)
			// Mark our selected left table cell
			username.SetTextColor(tcell.ColorRed)
			// Set our selected username
			inboxSelectedUsername = username.Text
			// Load our messages for the user
			go loadMessages(inboxSelectedUsername, inboxMessageContainer)
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
		if user.Username != lUser {
			userListContainer.SetCell(count, 0, tview.NewTableCell(user.Username))
			count++
		}
	}
	app.Draw()
}

func loadMessages(username string, container *tview.TextView) {
	inboxMessageContainerLock.Lock()
	defer inboxMessageContainerLock.Unlock()
	// Get our messages
	messages, _ := dbMessagesGet(username, lUser)
	reverseAny(messages)
	var result string
	for _, message := range messages {
		var fmsg string
		var color string
		var clearText string
		// Attempt to decrypt message
		// Note: This is not very efficient, the message may not be one of our own and hence
		// will fail increasing load time
		if message.To == lUser {
			s, err := decryptMessage(message.ToMessage)
			if err != nil {
				return
			}
			clearText = s
		} else {
			s, err := decryptMessage(message.FromMessage)
			if err != nil {
				return
			}
			clearText = s
		}
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
			fmsg += " <[darkmagenta]" + message.From + color +  "> [lightgray]"
		} else {
			fmsg += " <[darkcyan]" + message.From + color + "> [lightgray]"
		}
		// Set our message
		fmsg += clearText
		// Finalize our string
		result += fmsg + tablewriter.NEWLINE
	}
	// Clear our messages
	container.Clear()
	// Set our new message
	container.SetText(result)
	container.ScrollToEnd()
}

func inboxRetryFailedMessages(username string) {
	// reset counter
	inboxFailedMessageCount = 0
	// Get our messages
	messages, _ := dbMessagesGet(username, lUser)
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

func submitMessage(toUser string, message string) {
	var pgp = gopenpgp.GetGopenPGP()

	// Make sure we have this user in our local DB to encrypt
	usrObj, err := dbUserGet(toUser)
	if err != nil {
		showError(ClientError{
			Message:  "Unable to submit message to user as it does not exist",
			Button:   "Continue",
		})
		return
	}

	// Encrypt message using our recipients public key
	toKeyRing, err := gopenpgp.ReadArmoredKeyRing(strings.NewReader(usrObj.PubKey))
	if err != nil {
		panic(err)
	}
	encToMsg, err := pgp.EncryptMessage(
		message,
		toKeyRing,
		nil,
		"",
		false,
	)
	if err != nil {
		panic(err)
	}

	// Encrypt our message using our details
	fromKeyRing, err := gopenpgp.ReadArmoredKeyRing(strings.NewReader(privKey))
	if err != nil {
		panic(err)
	}
	encFromMsg, err := pgp.EncryptMessage(
		message,
		fromKeyRing,
		nil,
		"",
		false,
	)
	if err != nil {
		panic(err)
	}

	// Define our message object
	msgObj := &models.Message{
		FromMessage: encFromMsg,
		ToMessage: encToMsg,
		To:      toUser,
		From: lUser,
		Date:    time.Now(),
		Success: true,
	}

	// Add our message to our local DB
	id, err := dbMessageAdd(msgObj, msgObj.To)
	if err != nil {
		panic(err)
	}

	// Update our object with our db ID
	msgObj.ID = id

	client.cmdMsgTo(msgObj)
}