package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/mrdolev/sieve-go/internal/collector"
	ip_stats "github.com/mrdolev/sieve-go/internal/ip_stats"
	router "github.com/mrdolev/sieve-go/internal/router"
	"github.com/mrdolev/sieve-go/internal/storage"
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

	sharedStorage := storage.NewRedisClient(redisClient)
	ipStatsRepo := ip_stats.NewIPStatsRepo(sharedStorage)

	var collector collector.CollectorI = collector.NewCollector(
		ipStatsRepo,
		1000,
		5*time.Second,
	)

	defer collector.Close()

	var ipStatsService ip_stats.IPStatsServiceI = ip_stats.NewIPStatsService(collector)
	var handler ip_stats.IPStatsHandlerI = ip_stats.NewHandler(ipStatsService)

	serverPort := getEnvInt("SERVER_PORT", 8080)
	serverPath := getEnv("SERVER_PATH", "/")

	mux := http.NewServeMux()

	r := router.NewRouter(handler, mux)
	r.Collect(serverPath)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", serverPort),
		Handler: mux,
	}

	log.Println("server is started")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalln("error to start web server")
	}
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
