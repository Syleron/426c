package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/syleron/426c/common/database"
	"github.com/syleron/426c/common/security"
	"log"
	"strconv"
	"time"
)

var (
	app    = tview.NewApplication()
	pages  = tview.NewPages()
	layout *tview.Flex
	client *Client
	userList *UserList
	db     *database.Database
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

	keysHelp := "[gray](ESC) Exit/Back (TAB) Navigate (ENTER) Select [white]"


	// Do it first
	fmt.Fprintf(foot, keysHelp + " v" + Version())

	// Then update every 2 seconds
	go doEvery(2*time.Second, func() error {
		app.QueueUpdateDraw(func() {
			if _, err := client.Cache.Get("username"); err != nil {
				return
			}
			msgCost, err := client.Cache.Get("msgCost")
			if err != nil {
				log.Fatal(err)
			}
			blocks, err := client.Cache.Get("blocks")
			if err != nil {
				log.Fatal(err)
			}
			foot.Clear()
			fmt.Fprintf(foot, keysHelp+" [_] "+strconv.Itoa(blocks.(int))+" / "+strconv.Itoa(msgCost.(int)))
		})
		return nil
	})

	return foot
}

func main() {
	var err error
	mainCheckKeys()
	// Load our database
	db, err = database.New("426c")
	if err != nil {
		panic(err)
	}
	db.CreateBucket("messages")
	db.CreateBucket("users")
	// Setup our socket client
	client, err = setupClient("proteus.426c.net:9000")
	// Setup our user list
	userList = NewUserList(userListContainer)
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
