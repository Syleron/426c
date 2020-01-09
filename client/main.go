package main

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/syleron/426c/common/security"
	"strconv"
	"time"
)

var (
	app = tview.NewApplication()
	pages = tview.NewPages()
	layout *tview.Flex
	//user = &User{}
	//sockets *client.Client
	client *Client
)

func header() *tview.TextView {
	head := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false)

	fmt.Fprintf(head, `[:bu]░ [yellow]426c [white]Network `)

	return head
}

func footer() *tview.TextView {
	foot := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignRight).
		SetWrap(false)

	// Do it first
	fmt.Fprintf(foot, " [_] " + strconv.Itoa(getBlocks()) + " ")
	// Then update every 2 seconds
	go doEvery(2 * time.Second, func() error {
		foot.Clear()
		fmt.Fprintf(foot, " [_] " + strconv.Itoa(getBlocks()) + " ")
		app.Draw()
		return nil
	})

	return foot
}

func main() {
	// Generate our connection keys
	if err := security.GenerateKeys("127.0.0.1"); err != nil {
		panic(err)
	}
	// Setup our socket client
	var err error
	client, err = setupClient()
	// Defer our client close
	defer client.Close()
	// Create the main layout
	layout = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header(), 1, 1, false).
		AddItem(pages, 0, 1, true).
		AddItem(footer(), 1, 1, false)
	// Load our pages
	LoadPages()
	// Input
	InputHandlers()
	if err != nil {
		pages.SwitchToPage("unavailable")
	}
	// Start our main app loop
	if err := app.SetRoot(layout, true).Run(); err != nil {
		panic(err)
	}
}