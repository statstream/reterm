package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/caarlos0/env"
	"github.com/gdamore/tcell/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rivo/tview"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

var envConfig EnvironmentConfig = EnvironmentConfig{}
var pages = tview.NewPages()
var app = tview.NewApplication()
var keys = tview.NewList().ShowSecondaryText(false)
var valuesTextView = tview.NewTextView()
var flex = tview.NewFlex()
var searchInput *tview.InputField
var text = tview.NewTextView().
	SetTextColor(tcell.ColorGreen).
	SetText("\n(r) to refresh\n(d) to delete\n(q) to quit\n(h) for help\n(/) go to search bar")
var ctx = context.Background()
var redisClient *redis.Client

func main() {

	// envConfig := EnvironmentConfig{}
	err := env.Parse(&envConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse environment variables")
	}

	//Create a new instance of lumberjack.Logger to handle log file rotation
	logFile := &lumberjack.Logger{
		Filename:   envConfig.FileName,
		MaxSize:    envConfig.MaxSize,
		MaxBackups: envConfig.MaxBackups,
		MaxAge:     envConfig.MaxAge,
	}
	logger := zerolog.New(logFile).With().Timestamp().Caller().Logger()
	log.Logger = logger

	// Parsing Redis URL
	redisURL, err := redis.ParseURL(envConfig.RedisURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse Redis URL")
	}
	// Create a new Redis client
	redisClient = redis.NewClient(redisURL)

	// Input field for search keys
	searchInput = tview.NewInputField().SetLabel("Search Key: ").SetFieldBackgroundColor(tcell.ColorDefault).
		SetFieldTextColor(tcell.ColorWhite).SetFieldWidth(50).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			searchKeys(searchInput.GetText())
		}
	})

	// Fetch keys from Redis using the KEYS command
	keyItems, err := redisClient.Keys(ctx, "*").Result()
	if err != nil {
		panic(err)
	}

	// Add items to the "keys" list.
	for idx, key := range keyItems {
		keys.AddItem(fmt.Sprintf("%d. %s", idx+1, key), "", 0, nil)
	}

	// Scroll to a specific key in the "keys" section
	scrollToKeySection := func(key string) {
		for idx, keyItem := range keyItems {
			if strings.Contains(keyItem, key) {
				keys.SetCurrentItem(idx)
				break
			}
		}
	}
	// Houring
	keys.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		actualKey := strings.TrimSpace(strings.TrimPrefix(mainText, fmt.Sprintf("%d.", index+1)))
		hoverOverKey(actualKey)
	})
	keys.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		actualKey := strings.TrimSpace(strings.TrimPrefix(mainText, fmt.Sprintf("%d.", index+1)))
		hoverOverKey(actualKey)
	})

	// Set the selected function for the "keys" list
	keys.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		// Remove the serial number from the mainText to get the actual key
		actualKey := strings.TrimSpace(strings.TrimPrefix(mainText, fmt.Sprintf("%d.", index+1)))
		// houring the key value pair
		hoverOverKey(actualKey)
	})

	valuesTextView.SetBorder(true).SetTitle("Values")

	// Set up the layout with List, TextView, and Quit/Refresh buttons
	flex.SetDirection(tview.FlexRow).
		AddItem(searchInput, 2, 1, false).
		AddItem(tview.NewFlex().
			AddItem(keys, 0, 1, true).
			AddItem(valuesTextView, 0, 1, false), 0, 6, false).
		AddItem(text, 0, 1, false)

	// Set the input capture for handling 'q' (quit) and 'r' (refresh)
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if searchInput.HasFocus() {
			if event.Rune() == '/' {
				app.SetFocus(searchInput)
			} else if event.Key() == tcell.KeyEscape {
				app.SetFocus(keys)      // Move focus back to the "keys" list
				searchInput.SetText("") // Clear the search input field
				searchKeys("")
			} else if event.Key() == tcell.KeyEnter {
				searchText := searchInput.GetText()
				if searchText != "" {
					searchKeys(searchText)
					app.SetFocus(keys)
					scrollToKeySection(searchText)
				}
			}
		} else {
			// When search input is not in focus, handle 'q' (quit), 'r' (refresh), and 'd' (delete).
			if event.Key() == tcell.KeyDown {
				keys.SetCurrentItem(keys.GetCurrentItem() + 1)
			} else if event.Key() == tcell.KeyUp {
				keys.SetCurrentItem(keys.GetCurrentItem() - 1)
			} else if event.Rune() == 'q' {
				app.Stop()
			} else if event.Rune() == 'r' {
				refreshData()
			} else if event.Rune() == 'd' {
				selectedIndex := keys.GetCurrentItem()
				selectedText, _ := keys.GetItemText(selectedIndex)
				selectedKey := strings.TrimSpace(strings.TrimPrefix(selectedText, fmt.Sprintf("%d.", selectedIndex+1)))
				confirmDeleteModal(selectedIndex, selectedKey)
			} else if event.Rune() == 'h' {
				help()
			} else if event.Rune() == '/' {
				app.SetFocus(searchInput)
			}
		}
		return event
	})
	// Add the layout to the pages and run the app
	pages.AddPage("Menu", flex, true, true)
	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		log.Fatal().Err(err).Msg("Error running app")
	}

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		log.Fatal().Err(err).Msg("Error running app")

	}
	// Close the Redis client
	err = redisClient.Close()
	if err != nil {
		log.Error().Err(err).Msg("Error closing Redis client")
	}

}
