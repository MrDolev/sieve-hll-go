package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/mrdolev/sieve-go/internal/server"
	"github.com/redis/go-redis/v9"
)

func main() {
	addr := getEnv("REDIS_ADDR", "localhost:6379")
	password := getEnv("REDIS_PASSWORD", "")
	db := getEnvInt("REDIS_DB", 0)
	proto := getEnvInt("REDIS_PROTOCOL", 2)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		Protocol: proto,
	})

	ctx := context.Background()

	err := redisClient.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		log.Fatalf("panic to set %s", err)
	}

	result, err := redisClient.Get(ctx, "foo").Result()
	if err != nil {
		log.Fatalf("error %s", err)
	}
	log.Printf("query result: %s", result)

	serverPort := getEnvInt("SERVER_PORT", 8080)
	serverPath := getEnv("SERVER_PATH", "/")

	serverMux := server.NewServerMux(serverPort, serverPath)
	serverMux.Serve()
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err == nil && value >= 0 {
		return value
	}
	return fallback
}
