package main

import "github.com/rivo/tview"

type Page func() (id string, content tview.Primitive)

var (
	pageItems = []Page{
		SplashPage,
		LoginPage,
		RegisterPage,
		RegisterWarningPage,
		InboxPage,
		UnavailablePage,
	}
)

func LoadPages() {
	for i, p := range pageItems {
		title, primitive := p()
		pages.AddPage(title, primitive, true, i == 0)
	}
}
