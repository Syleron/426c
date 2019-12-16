package main

import "github.com/rivo/tview"

type ClientError struct {
	Message string
	Button string
	Continue func()
}

func showError(mError ClientError) {
	modal := tview.NewModal().
		SetText(mError.Message).
		AddButtons([]string{mError.Button}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(pages, true)
			mError.Continue()
		})
	app.SetRoot(modal, true)
}