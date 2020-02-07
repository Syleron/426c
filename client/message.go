package main

import (
	gopenpgp "github.com/ProtonMail/gopenpgp/crypto"
	"github.com/olekukonko/tablewriter"
	"github.com/rivo/tview"
	"github.com/syleron/426c/common/models"
	"strings"
	"time"
)

func messageLoad(username string, container *tview.TextView) {
	// Get our messages
	messages, _ := dbMessagesGet(username, client.Username)
	reverseAny(messages)
	var result string
	for _, message := range messages {
		var fmsg string
		var color string
		var clearText string

		// Attempt to decrypt message
		// Note: This is not very efficient, the message may not be one of our own and hence
		// will fail increasing load time
		if message.To == client.Username {
			s, err := messageDecrypt(message.ToMessage)
			if err != nil {
				return
			}
			clearText = s
		} else {
			s, err := messageDecrypt(message.FromMessage)
			if err != nil {
				return
			}
			clearText = s
		}

		// Start structuring our message
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
		if message.To == client.Username {
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

func messageSubmit(toUser string, message string) {
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
		From: client.Username,
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

func messageDecrypt(message string) (string, error) {
	var pgp = gopenpgp.GetGopenPGP()
	clearText, err := pgp.DecryptMessageStringKey(message, privKey, pHash)
	if err != nil {
		return "", err
	}
	return clearText, nil
}
