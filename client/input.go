package main

import "github.com/gdamore/tcell"

func InputHandlers() {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlS:
			//if loggedIn() {
			//	pages.SwitchToPage("search")
			//}
		case tcell.KeyCtrlQ:
			app.Stop()
		case tcell.KeyCtrlC:
			//if loggedIn() {
			//	pages.SwitchToPage("chat")
			//}
			return nil
		}
		return event
	})
}