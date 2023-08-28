package main

import (
	"github.com/rivo/tview"
)

// For help
func help() {
	helpText := "Press 'd' key to delete the selected key from Redis.\n" +
		"Press 'q' key to quit the terminal.\n" +
		"Press 'r' key to refresh the keys in the terminal.\n" +
		"Press '/' to go to search bar.\n " +
		"Press  'ESC' to go back to key list from search bar."

	Model := tview.NewModal().
		SetText("Help").
		SetText(helpText).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetFocus(flex)
			pages.HidePage("Help")
		})

	pages.AddPage("Help", Model, true, true)
	app.SetFocus(Model)
}
