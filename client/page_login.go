package main

import (
	"github.com/rivo/tview"
)

func LoginPage() (id string, content tview.Primitive) {
	var username, password string
	form := tview.NewForm().
		AddInputField("Username:", "", 20, nil, func(text string) {
			username = text
		}).
		AddPasswordField("Password:", "", 20, '*', func(text string) {
			password = text
		}).
		AddButton("Login", func() {
			if username == "" || password == "" {
				showError(ClientError{
					Message:  "Please enter a username and password",
					Button:   "Continue",
					Continue: func() {
						pages.SwitchToPage("login")
					},
				})
				return
			}
		// Submit our registration to the server
		if err := client.msgLogin(username, password); err != nil {
			panic(err)
		}
	}).
		AddButton("Register", func() {
			pages.SwitchToPage("register warning")
		}).
		AddButton("Exit", func() {
			app.Stop()
		})
	form.SetBorder(true).SetTitle("426c Login")

	return "login", Center(40, 10, form)
}
