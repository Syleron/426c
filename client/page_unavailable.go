package main

import (
	"github.com/rivo/tview"
)

// TODO: Convert to modal
func UnavailablePage() (id string, content tview.Primitive) {
	modal := tview.NewModal().
		SetText("Unable to communicate with the 426c network. Please try again later.").
		AddButtons([]string{"Exit"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Exit" {
				app.Stop()
			}
		})
	return "unavailable", modal
}
