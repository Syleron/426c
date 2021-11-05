package main

type ClientError struct {
	Message  string
	Button   string
	Continue func()
}

func showError(mError ClientError) {
	//modal := tview.NewModal().
	//	SetText(mError.Message).
	//	AddButtons([]string{mError.Button}).
	//	SetDoneFunc(func(buttonIndex int, buttonLabel string) {
	//		app.SetRoot(layout, true)
	//		if mError.Continue != nil {
	//			mError.Continue()
	//		}
	//	})
	//app.QueueUpdateDraw(func() {
	//	app.SetRoot(modal, true)
	//})
}
