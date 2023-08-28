package main

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

// Searching Key
func searchKeys(query string) {
	keyItems, err := redisClient.Keys(ctx, "*").Result()
	if err != nil {
		log.Error().Err(err).Msg("Error for searching keys")
	}

	filteredKeys := make([]string, 0)

	for _, key := range keyItems {
		if strings.Contains(key, query) {
			filteredKeys = append(filteredKeys, key)
		}
	}
	keys.Clear()
	for idx, key := range filteredKeys {
		keys.AddItem(fmt.Sprintf("%d. %s", idx+1, key), "", 0, nil)

	}
}
