package main

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/syleron/426c/common/security"
)

var (
	app = tview.NewApplication()
	pages = tview.NewPages()
	//user = &User{}
	//sockets *client.Client
	client *Client
)

func header() *tview.TextView {
	head := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false)

	//if user.Token != "" {
	fmt.Fprintf(head, `░ 426c Network`)
	//}

	return head
}

func footer() *tview.TextView {
	foot := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignRight).
		SetWrap(false)

	fmt.Fprintf(foot, "Connected <")

	return foot
}

func main() {
	// Generate our connection keys
	if err := security.GenerateKeys("127.0.0.1"); err != nil {
		panic(err)
	}
	// Setup our socket client
	client = setupClient()
	// Defer our client close
	defer client.Close()
	// Put our handlers into a go rutine
	go client.connectionHandler()
	// Create the main layout
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header(), 1, 1, false).
		AddItem(pages, 0, 1, true).
		AddItem(footer(), 1, 1, false)
	// Load our pages
	LoadPages()
	// Input
	InputHandlers()
	// Start our main app loop
	if err := app.SetRoot(layout, true).Run(); err != nil {
		panic(err)
	}
}