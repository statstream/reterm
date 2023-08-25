package main

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

// Delete selected key
func deleteSelectedKey(selectedIndex int, selectedKey string) {
	_, err := redisClient.Del(ctx, selectedKey).Result()
	if err != nil {
		valuesTextView.SetText("Error deleting key from Redis")
		log.Error().Err(err).Str("key", selectedKey).Msg("Error deleting key from Redis")

	} else {
		log.Info().Str("key", selectedKey).Msg("Deleted key from Redis")
		valuesTextView.SetText(fmt.Sprintf("Key '%s' deleted successfully", selectedKey))
		keys.RemoveItem(selectedIndex)
	}
}
func confirmDeleteModal(selectedIndex int, selectedKey string) {
	modal := tview.NewModal().
		SetText("Are you sure you want to delete the selected key?").
		AddButtons([]string{" Yes ", " No "}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == " Yes " {
				deleteSelectedKey(selectedIndex, selectedKey)
				refreshData()
			}
			app.SetFocus(keys)
			pages.RemovePage("modal")
		})

	pages.AddPage("modal", modal, true, true)
	app.SetFocus(modal)
}
