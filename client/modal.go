package main

import "github.com/rivo/tview"

type ClientModal struct {
	Message      string
	SubmitButton string
	CancelButton string
	Continue     func()
	Cancel       func()
}

func showModal(mModal ClientModal) {
	var buttons []string
	// Add the submit button if the text is set
	if mModal.SubmitButton != "" {
		buttons = append(buttons, mModal.SubmitButton)
	}
	// Add the cancel button if the text is set
	if mModal.CancelButton != "" {
		buttons = append(buttons, mModal.CancelButton)
	}
	modal := tview.NewModal().
		SetText(mModal.Message).
		AddButtons(buttons)
	// Handle button selection
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		app.SetRoot(layout, true)
		switch buttonIndex {
		case 1:
			if mModal.Cancel != nil {
				mModal.Cancel()
			}
		default:
			if mModal.Continue != nil {
				mModal.Continue()
			}
		}
	})
	app.QueueUpdateDraw(func() {
		app.SetRoot(modal, true)
	})
}
