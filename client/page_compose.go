package main

import (
	gopenpgp "github.com/ProtonMail/gopenpgp/crypto"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/syleron/426c/common/models"
	"github.com/syleron/426c/common/packet"
	"github.com/syleron/426c/common/utils"
	"github.com/syleron/femto"
	"strings"
	"time"
)

var (
	composeMessageContainer *tview.TextView
	composeMessageField     *tview.InputField
	composeToField          *tview.InputField
)

func ComposePage() (id string, content tview.Primitive) {
	grid := tview.NewGrid().
		SetRows(1).
		SetColumns(30, 0).
		SetBorders(false).
		SetGap(0, 2)

	userGrid := tview.NewFlex()
	chatGrid := tview.NewFlex()

	chatGrid.SetBorder(true)
	chatGrid.SetBorderPadding(1, 1, 1, 1)
	chatGrid.SetTitle(" Compose New Message ")

	messageContainer = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	messageContainer.SetScrollable(true)

	buffer := femto.NewBufferFromString("\n\n------\nThis is an encrypted message sent via 426c", "")
	inputField := femto.NewView(buffer)

	toInputField := tview.NewInputField().
		SetPlaceholder("Enter username")

	// Cancel button
	cancelButton := tview.NewButton("Cancel").SetSelectedFunc(func() {
		pages.SwitchToPage("inbox")
	})
	cancelButton.SetBorder(true).SetRect(0, 0, 0, 1)

	// Send button
	sendButton := tview.NewButton("Send Message").SetSelectedFunc(func() {
		submitMessage(toInputField.GetText(), buffer.String())
	})
	sendButton.SetBorder(true).SetRect(0, 0, 0, 1)

	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlS:
			//saveBuffer(buffer, path)
			return nil
		case tcell.KeyCtrlQ:
			//app.Stop()
			return nil
		}
		return event
	})

	toInputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTAB:
			// Check if user exists and get public key details
			_, err := client.Send(packet.CMD_USER, utils.MarshalResponse(&models.UserRequestModel{
				Username: toInputField.GetText(),
			}))
			if err != nil {
				panic(err)
			}
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
			// Redraw our contacts list
			drawContactsList()
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
		Date:    time.Time{},
		Success: false,
	}

	// Add our message to our local DB
	id, err := dbMessageAdd(msgObj)
	if err != nil {
		panic(err)
	}

	// Update our object with our db ID
	msgObj.ID = id

	// Add our message to our message queue to send/process
	client.MQ.Add(msgObj)

	// Process our message queue
	go client.MQ.Process()

	// Switch back to our inbox
	pages.SwitchToPage("inbox")
}
