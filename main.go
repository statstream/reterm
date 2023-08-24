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

var text = tview.NewTextView().
	SetTextColor(tcell.ColorGreen).
	SetText("\n(q) to Quit\n(r) to Refresh")

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

	// Initialize a Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr:     envConfig.RedisAddress,
		Password: envConfig.RedisPassword,
		DB:       envConfig.RedisDB,
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

	// Set the selected function for the "keys" list
	keys.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		// Remove the serial number from the mainText to get the actual key
		actualKey := strings.TrimSpace(strings.TrimPrefix(mainText, fmt.Sprintf("%d.", index+1)))

		valueType, err := redisClient.Type(ctx, actualKey).Result()
		if err != nil {
			valuesTextView.SetText("Error fetching key type from Redis")
			return
		}

		switch valueType {
		case "string":
			value, err := redisClient.Get(ctx, actualKey).Result()
			if err == nil {
				log.Info().Str("key", mainText).Msg("Fetched string value from Redis")
				valuesTextView.SetText(value)
			} else if err == redis.Nil {
				log.Info().Str("key", mainText).Msg("Key not found in Redis")
				valuesTextView.SetText("No value available")
			} else {
				log.Error().Err(err).Str("key", mainText).Msg("Error fetching value from Redis")
				valuesTextView.SetText("Error fetching value from Redis")
			}
		case "list":
			listValues, err := redisClient.LRange(ctx, actualKey, 0, -1).Result()
			if err != nil {
				log.Error().Err(err).Str("key", mainText).Msg("Error fetching list values from Redis")
				valuesTextView.SetText("Error fetching list values from Redis")
				return
			}
			log.Info().Str("key", mainText).Msg("Fetched list values from Redis")

			valuesTextView.SetText(strings.Join(listValues, "\n")) // Display list values
		case "hash":
			hashValues, err := redisClient.HGetAll(ctx, actualKey).Result()
			if err != nil {
				log.Error().Err(err).Str("key", mainText).Msg("Error fetching hash values from Redis")
				valuesTextView.SetText("Error fetching hash values from Redis")
				return
			}
			log.Info().Str("key", mainText).Msg("Fetched hash values from Redis")

			var hashValueStrings []string
			for field, value := range hashValues {
				hashValueStrings = append(hashValueStrings, fmt.Sprintf("%s: %s", field, value))
			}
			valuesTextView.SetText(strings.Join(hashValueStrings, "\n")) // Display hash values
		default:
			log.Warn().Str("key", mainText).Msg("Unknown value type")
			valuesTextView.SetText("Unknown value type")
		}
	})

	valuesTextView.SetBorder(true).SetTitle("Value Of Key")

	// Set up the layout with List, TextView, and Quit/Refresh buttons
	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(keys, 0, 1, true).
			AddItem(valuesTextView, 0, 1, false), 0, 6, false).
		AddItem(text, 0, 1, false)

	// Set the input capture for handling 'q' (quit) and 'r' (refresh)
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			app.Stop()
		} else if event.Rune() == 'r' {
			refreshData()
		}
		return event
	})
	// Add the layout to the pages and run the app
	pages.AddPage("Menu", flex, true, true)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		log.Fatal().Err(err).Msg("Error running app")

	}
	// Close the Redis client
	err = redisClient.Close()
	if err != nil {
		log.Error().Err(err).Msg("Error closing Redis client")
	}
}
