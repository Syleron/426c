package main

import "github.com/rivo/tview"

var (
	pageItems = []Page{
		SplashPage,
		LoginPage,
		RegisterPage,
		RegisterWarningPage,
		InboxPage,
		SearchPage,
		ComposePage,
		UnavailablePage,
	}
)

type Page func() (id string, content tview.Primitive)

func LoadPages() {
	for i, p := range pageItems {
		title, primitive := p()
		pages.AddPage(title, primitive, true, i == 0)
	}
}