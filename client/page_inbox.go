package main

import "github.com/rivo/tview"

func InboxPage() (id string, content tview.Primitive) {
	var username, password string
	form := tview.NewForm().
		AddInputField("Username:", "", 20, nil, func(text string) {
			username = text
		}).
		AddPasswordField("Password:", "", 20, '*', func(text string) {
			password = text
		}).
		AddButton("Login", func() {
		}).
		AddButton("Register", func() {
			pages.SwitchToPage("register warning")
		}).
		AddButton("Exit", func() {
			app.Stop()
		})
	form.SetBorder(true).SetTitle("426c Login")

	return "inbox", Center(40, 10, form)
}
