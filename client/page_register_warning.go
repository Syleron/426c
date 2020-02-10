package main

import "github.com/rivo/tview"

// TODO: Convert to modal
func RegisterWarningPage() (id string, content tview.Primitive) {
	modal := tview.NewModal().
		SetText("Your encryption is only as good as your password. You cannot recover or change your password.").
		AddButtons([]string{"I Understand"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "I Understand" {
				pages.SwitchToPage("register")
			}
		})
	return "register warning", modal
}
