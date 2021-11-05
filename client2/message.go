package main

import (
	gopenpgp "github.com/ProtonMail/gopenpgp/crypto"
	"github.com/olekukonko/tablewriter"
	"github.com/rivo/tview"
	"github.com/syleron/426c/common/models"
	"log"
	"strings"
	"time"
)


// TODO: This needs to change so that there is some sort of cache in memory. As it's currently slowwww
func messageLoad(username string, container *tview.TextView) {
	clientUsername, err := client.Cache.Get("username")
	if err != nil {
		log.Fatal(err)
	}
	// Get our messages
	messages, _ := dbMessagesGet(username, clientUsername.(string))
	reverseAny(messages)
	var result string
	for _, message := range messages {
		var fmsg string
		var color string
		var clearText string

		// Attempt to decrypt message
		// Note: This is not very efficient, the message may not be one of our own and hence
		// will fail increasing load time
		if message.To == clientUsername.(string) {
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
		if message.To == clientUsername.(string) {
			fmsg += " <[darkmagenta]" + message.From + color +  "> [lightgray]"
		} else {
			fmsg += " <[darkcyan]" + message.From + color + "> [lightgray]"
		}
		// Set our message
		fmsg += clearText

		// Finalize our string
		result += fmsg + tablewriter.NEWLINE
	}
	// Set our new message
	//app.QueueUpdateDraw(func() {
	//	container.SetText(result)
	//	container.ScrollToEnd()
	//})
}

func messageSubmit(toUser string, message string) {
	var pgp = gopenpgp.GetGopenPGP()

	clientUsername, err := client.Cache.Get("username")
	if err != nil {
		log.Fatal(err)
	}
	pKey, err := client.Cache.Get("pKey")
	if err != nil {
		log.Fatal(err)
	}

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
		log.Fatal(err)
	}

	// Encrypt our message using our details
	fromKeyRing, err := gopenpgp.ReadArmoredKeyRing(strings.NewReader(pKey.(string)))
	if err != nil {
		log.Fatal(err)
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
		From: clientUsername.(string),
		Date:    time.Now(),
		Success: true,
	}

	// Add our message to our local DB
	id, err := dbMessageAdd(msgObj, msgObj.To)
	if err != nil {
		log.Fatal(err)
	}

	// Update our object with our db ID
	msgObj.ID = id

	client.cmdMsgTo(msgObj)
}

func messageDecrypt(message string) (string, error) {
	var pgp = gopenpgp.GetGopenPGP()

	passHash, err := client.Cache.Get("passHash")
	if err != nil {
		log.Fatal(err)
	}
	pKey, err := client.Cache.Get("pKey")
	if err != nil {
		log.Fatal(err)
	}

	clearText, err := pgp.DecryptMessageStringKey(message, pKey.(string), passHash.(string))
	if err != nil {
		return "", err
	}
	return clearText, nil
}
