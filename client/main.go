package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/syleron/426c/common/database"
	"github.com/syleron/426c/common/security"
	"github.com/syleron/426c/common/utils"
	"strconv"
	"time"
)

var (
	app    = tview.NewApplication()
	pages  = tview.NewPages()
	layout *tview.Flex
	client *Client
	db     *database.Database

	// User password hash used to decrypt w/ private key.
	// TODO: Not sure if this should be done or not. Seems iffy.
	pHash string
	// Logged in user private key
	privKey string
	// Total blocks available
	blocks int
)

func header() *tview.TextView {
	head := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false)

	fmt.Fprintf(head, `░ [yellow]426c [gray]Network `)

	return head
}

func footer() *tview.TextView {
	foot := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetTextAlign(tview.AlignRight).
		SetWrap(false)

	foot.SetBackgroundColor(tcell.NewRGBColor(48, 48, 48))

	// Do it first
	fmt.Fprintf(foot, " [_] "+strconv.Itoa(blocks)+" ")
	// Then update every 2 seconds
	go doEvery(2*time.Second, func() error {
		foot.Clear()
		fmt.Fprintf(foot, " [_] "+strconv.Itoa(blocks)+" ")
		app.Draw()
		return nil
	})

	return foot
}

func main() {
	var err error
	mainCheckKeys()
	mainLoadPrivateKey()
	// Load our database
	db, err = database.New("426c")
	if err != nil {
		panic(err)
	}
	db.CreateBucket("messages")
	db.CreateBucket("users")
	// Setup our socket client
	client, err = setupClient("proteus.426c.net:9000")
	// Create the main layout
	layout = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header(), 1, 1, false).
		AddItem(pages, 0, 1, true).
		AddItem(footer(), 1, 1, false)
	// Load our pages
	LoadPages()
	if err != nil {
		pages.SwitchToPage("unavailable")
	}
	// Start our main app loop
	if err := app.SetRoot(layout, true).Run(); err != nil {
		panic(err)
	}
}

func mainCheckKeys() {
	// Generate our connection keys
	if err := security.GenerateKeys("127.0.0.1"); err != nil {
		panic(err)
	}
}

func mainLoadPrivateKey() {
	// Load our key into memory
	b, err := utils.LoadFile("key.pem")
	if err != nil {
		panic(err)
	}
	// Set our private key
	privKey = string(b)
}
