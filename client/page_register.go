package main

import "github.com/rivo/tview"

func RegisterPage() (id string, content tview.Primitive) {
	var username, password, passwordAgain string
	form := tview.NewForm().
		AddInputField("Username:", "", 20, nil, func(text string) {
			username = text
		}).
		AddPasswordField("Password:", "", 20, '*', func(text string) {
			password = text
		}).
		AddPasswordField("Password Again:", "", 20, '*', func(text string) {
			passwordAgain = text
		}).
		AddButton("Register", func() {
			// Make sure our passwords match!
			if password != passwordAgain {
				showError(ClientError{
					Message:  "Passwords do not match",
					Button:   "Continue",
					Continue: func() {
						pages.SwitchToPage("register")
					},
				})
				return
			}
			// Submit our registration to the server
			client.msgRegister(username, password)
			pages.SwitchToPage("login")
		}).
		AddButton("Back", func() {
			pages.SwitchToPage("login")
		})
	form.SetBorder(true).SetTitle("426c Register")

	return "register", Center(40, 11, form)
}
