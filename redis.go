package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Client *redis.Client
}

func initRedisClient() (RedisConfig, error) {

	redisAddress := os.Getenv("REDIS_ADDRESS")
	if redisAddress == "" {
		return RedisConfig{}, fmt.Errorf("REDIS_ADDRESS environment variable is not set")
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := os.Getenv("REDIS_DB")
	if redisDB == "" {
		return RedisConfig{}, fmt.Errorf("REDIS_DB environment variable is not set")
	}

	db, err := strconv.Atoi(redisDB)
	if err != nil {
		return RedisConfig{}, fmt.Errorf("error converting REDIS_DB to int: %v", err)
	}

	redisClient := RedisConfig{
		Client: redis.NewClient(&redis.Options{
			Addr:     redisAddress,
			Password: redisPassword,
			DB:       db,
		}),
	}

	_, err = redisClient.Client.Ping(context.Background()).Result()
	if err != nil {
		return RedisConfig{}, fmt.Errorf("error connecting to Redis: %v", err)
	}
	fmt.Println("Connected to Redis:", redisClient.Client.Options().Addr)
	return redisClient, nil
}

func (redisClient RedisConfig) Close() error {
	err := redisClient.Client.Close()
	if err != nil {
		return fmt.Errorf("error closing Redis client: %v", err)
	}
	return nil
}

func (redisClient RedisConfig) FlushAll() error {
	err := redisClient.Client.FlushAll(context.Background()).Err()
	if err != nil {
		return fmt.Errorf("error flushing all keys: %v", err)
	}
	return nil
}

func (redisClient RedisConfig) SaveHTTPResponse(key string, response []byte) error {
	err := redisClient.Client.Set(context.Background(), key, response, 1*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("error saving HTTP response for key %s: %v", key, err)
	}
	return nil
}

func (redisClient RedisConfig) GetHTTPResponse(key string) ([]byte, error) {
	val, err := redisClient.Client.Get(context.Background(), key).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting HTTP response for key %s: %v", key, err)
	}
	return []byte(val), nil
}
func (redisClient RedisConfig) DeleteHTTPResponse(key string) error {
	err := redisClient.Client.Del(context.Background(), key).Err()
	if err != nil {
		return fmt.Errorf("error deleting HTTP response for key %s: %v", key, err)
	}
	return nil
}
func (redisClient RedisConfig) ExistsHTTPResponse(key string) (bool, error) {
	val, err := redisClient.Client.Exists(context.Background(), key).Result()
	if err != nil {
		return false, fmt.Errorf("error checking existence of HTTP response for key %s: %v", key, err)
	}
	return val > 0, nil
}
