package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
)

// Refreshing the key value pair
func refreshData() {
	keyItems, err := redisClient.Keys(ctx, "*").Result()

	if err != nil {
		log.Error().Err(err).Msg("Error fetching keys from redis while refreshing")
		return
	}
	keys.Clear()
	for idx, key := range keyItems {
		keys.AddItem(fmt.Sprintf("%d. %s", idx+1, key), "", 0, nil)
	}
	log.Info().Msg("Refreshed key list")

}
