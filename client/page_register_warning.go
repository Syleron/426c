package main

import "github.com/rivo/tview"

func RegisterWarningPage() (id string, content tview.Primitive) {
	modal := tview.NewModal().
		SetText("You're encryption is only as good as your password.").
		AddButtons([]string{"I Understand"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "I Understand" {
				pages.SwitchToPage("register")
			}
		})
	return "register warning", modal
}
