package main

import "github.com/rivo/tview"

// TODO: Replace this with a modal
func SplashPage() (id string, content tview.Primitive) {
	modal := tview.NewModal().
		SetText("This program is distributed in the hope that it will be legally useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.").
		AddButtons([]string{"Continue", "Exit"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Continue" {
				pages.SwitchToPage("login")
			} else if buttonLabel == "Exit" {
				app.Stop()
			}
		})

	return "splash", modal
}
