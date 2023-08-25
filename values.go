package main

import (
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// Display string values
func displayStringValue(key string) {
	value, err := redisClient.Get(ctx, key).Result()
	if err == nil {
		log.Info().Str("key", key).Msg("Fetched string value from Redis")
		valuesTextView.SetText(value)
	} else if err == redis.Nil {
		log.Info().Str("key", key).Msg("Key not found in Redis")
		valuesTextView.SetText("No value available")
	} else {
		log.Error().Err(err).Str("key", key).Msg("Error fetching value from Redis")
		valuesTextView.SetText("Error fetching value from Redis")
	}
}

// Display list values
func displayListValues(key string) {
	listValues, err := redisClient.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("Error fetching list values from Redis")
		valuesTextView.SetText("Error fetching list values from Redis")
		return
	}
	log.Info().Str("key", key).Msg("Fetched list values from Redis")

	valuesTextView.SetText(strings.Join(listValues, "\n"))
}

// / Display hash values
func displayHashValues(key string) {
	hashValues, err := redisClient.HGetAll(ctx, key).Result()
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("Error fetching hash values from Redis")
		valuesTextView.SetText("Error fetching hash values from Redis")
		return
	}
	log.Info().Str("key", key).Msg("Fetched hash values from Redis")

	var hashValueStrings []string
	for field, value := range hashValues {
		hashValueStrings = append(hashValueStrings, fmt.Sprintf("%s: %s", field, value))
	}

	valuesTextView.SetText(strings.Join(hashValueStrings, "\n"))
}

// Display set values
func displaySetValues(key string) {
	setValues, err := redisClient.SMembers(ctx, key).Result()
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("Error fetching set values from Redis")
		valuesTextView.SetText("Error fetching set values from Redis")
		return
	}
	log.Info().Str("key", key).Msg("Fetched set values from Redis")

	valuesTextView.SetText(strings.Join(setValues, "\n"))
}

// Display zset values
func displayZSetValues(key string) {
	zsetValues, err := redisClient.ZRangeWithScores(ctx, key, 0, -1).Result()
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("Error fetching zset values from Redis")
		valuesTextView.SetText("Error fetching zset values from Redis")
		return
	}
	log.Info().Str("key", key).Msg("Fetched zset values from Redis")

	var zsetStrings []string
	for _, z := range zsetValues {
		zsetStrings = append(zsetStrings, fmt.Sprintf("%s: %f", z.Member, z.Score))
	}

	valuesTextView.SetText(strings.Join(zsetStrings, "\n"))
}

// Display HyperLogLog value
func displayHyperLogLogValues(key string) {
	hllCount, err := redisClient.PFCount(ctx, key).Result()
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("Error fetching HyperLogLog value from Redis")
		valuesTextView.SetText("Error fetching HyperLogLog value from Redis")
		return
	}
	log.Info().Str("key", key).Msg("Fetched HyperLogLog value from Redis")

	valuesTextView.SetText(fmt.Sprintf("HyperLogLog count: %d", hllCount))
}

// Display Bitmap values
func displayBitmapValues(key string) {
	bitmap, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("Error fetching Bitmap value from Redis")
		valuesTextView.SetText("Error fetching Bitmap value from Redis")
		return
	}
	log.Info().Str("key", key).Msg("Fetched Bitmap value from Redis")

	valuesTextView.SetText(fmt.Sprintf("Bitmap value: %s", bitmap))
}
