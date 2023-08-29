package main

import (
	"github.com/rs/zerolog/log"
)

func hoverOverKey(key string) {

	valueType, err := redisClient.Type(ctx, key).Result()
	if err != nil {
		valuesTextView.SetText("Error fetching key type from Redis")
		log.Error().Err(err).Msg("Error fetching key type from Redis")
		return
	}

	// Display the value as per valueType
	switch valueType {
	case "string":
		displayStringValue(key)
	case "list":
		displayListValues(key)
	case "hash":
		displayHashValues(key, valueType)
	case "set":
		displaySetValues(key)
	case "zset":
		displayZSetValues(key)
	case "hyperloglog":
		displayHyperLogLogValues(key)
	case "bitmap":
		displayBitmapValues(key)
	default:
		log.Warn().Str("key", key).Msg("Unknown value type")
		valuesTextView.SetText("Unknown value type")
	}
}
