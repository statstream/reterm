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

	key1 := make([]string, 0)

	for _, key := range keyItems {
		if strings.Contains(key, query) {
			key1 = append(key1, key)
		}
	}
	keys.Clear()
	for idx, key := range key1 {
		keys.AddItem(fmt.Sprintf("%d. %s", idx+1, key), "", 0, nil)

	}
	if len(key1) > 0 {
		keys.SetCurrentItem(0)
	}
}
