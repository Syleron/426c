package main

import (
	g "github.com/AllenDang/giu"
	"github.com/labstack/gommon/log"
	"github.com/syleron/426c/common/database"
	"github.com/syleron/426c/common/security"
	"os"
)

var (
	//userList *UserList
	client *Client
	db     *database.Database

	showWarning bool
	showRegister bool
	showLogin bool
	showContactList bool

	showWindow2 bool
	checked     bool
)

func onShowWindow2() {
	showWindow2 = true
}

func onHideWindow2() {
	showWindow2 = false
}

func loop() {
	width, height := g.GetAvailableRegion()

	g.MainMenuBar().Layout(
		g.Menu("File").Layout(
			g.MenuItem("Exit").OnClick(func() {
				os.Exit(0)
			}),
		),
		g.Menu("Settings").Layout(),
		g.Menu("Help").Layout(),
	).Build()

	if showWarning {
		g.Window("WARNING").
			IsOpen(&showWarning).
			Flags(g.WindowFlagsNoCollapse | g.WindowFlagsNoResize).
			Pos((width / 2) - 40, (height / 2) - 40).
			Size(600, 110).
			Layout(
				g.Custom(func() {
					//width, _ := g.GetAvailableRegion()

					g.Row(
						g.Label("This program is distributed in the hope that it will be legally useful, but \nWITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY \nor FITNESS FOR A PARTICULAR PURPOSE."),
					).Build()
				}),
				g.Custom(func() { // centered buttons
					width, _ := g.GetAvailableRegion()
					buttonW := (width - 5)/2

					g.Row(
						g.Button("Continue").Size(buttonW, 25).OnClick(func() {
							showWarning = false
							showLogin = true
						}),
						g.Button("Exit").Size(buttonW, 25).OnClick(func() {
							os.Exit(0)
						}),
					).Build()
				}),
			)
	}

	if showLogin {
		var username string
		var password string
		g.Window("Login").IsOpen(&showWindow2).Flags(g.WindowFlagsNone).Pos((width / 2) - 40, (height / 2) - 40).Size(300, 110).Layout(
			g.Custom(func() {
				width, _ := g.GetAvailableRegion()

				g.Row(
					g.InputText(&username).Size(width),
				).Build()
				g.Row(
					g.InputText(&password).Size(width),
				).Build()
			}),
			g.Custom(func() { // centered buttons
				width, _ := g.GetAvailableRegion()
				buttonW := (width - 5)/2

				g.Row(
					g.Button("Login").Size(buttonW, 25).OnClick(nil),
					g.Button("Exit").Size(buttonW, 25).OnClick(func() {
						os.Exit(0)
					}),
				).Build()
			}),
		)
	}

	if showRegister {

	}

	if showContactList {
		g.Window("Contacts").IsOpen(&showWindow2).Flags(g.WindowFlagsNone).Pos(250, 30).Size(200, 100).Layout(
			g.Label("I'm a label in window 2"),
			g.Button("Hide me").OnClick(onHideWindow2),
		)
	}
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
	if err != nil {
		log.Fatal("Unable to connect to 426c network")
	}
	// Setup our user list

	wnd := g.NewMasterWindow("426c Network", 800, 600, g.MasterWindowFlagsNotResizable)
	showWarning = true
	wnd.Run(loop)
}

func mainCheckKeys() {
	// Generate our connection keys
	if err := security.GenerateKeys("127.0.0.1"); err != nil {
		panic(err)
	}
}
