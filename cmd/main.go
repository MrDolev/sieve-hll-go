package main

import (
	"fmt"
	"log"
	"net/http"
	"net/netip"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mrdolev/sieve-hll-go/internal/collector"
	ip_stats "github.com/mrdolev/sieve-hll-go/internal/ip_stats"
	netfilter "github.com/mrdolev/sieve-hll-go/internal/net_filter"
	router "github.com/mrdolev/sieve-hll-go/internal/router"
	"github.com/mrdolev/sieve-hll-go/internal/storage"
	"github.com/redis/go-redis/v9"
)

func main() {
	addr := getEnv("REDIS_ADDR", "localhost:6379")
	password := getEnv("REDIS_PASSWORD", "")
	db := getEnvInt("REDIS_DB", 0)
	proto := getEnvInt("REDIS_PROTOCOL", 2)
	cidr, err := getEnvCIDR("CIDR")
	if err != nil {
		log.Fatalf("error to parse configuration: %s", err)
	}
	trustedProxies, err := getTrustedProxies()
	if err != nil {
		log.Fatalf("error to parse configuration: %s", err)
	}
	if len(trustedProxies) == 0 {
		log.Println("trusted proxies: none")
	} else {
		log.Printf("trusted proxies: %v", trustedProxies)
	}
	log.Println("warning: TRUSTED_PROXIES is not enforced yet, X-Forwarded-For / X-Real-IP are still read from every client")

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

	var netFilter netfilter.NetFilterI = netfilter.NewNetFilter(cidr)
	var ipStatsService ip_stats.IPStatsServiceI = ip_stats.NewIPStatsService(collector, netFilter)
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
		log.Fatalf("error to start web server: %s", err)
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

func getEnvList(key string, fallback []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return strings.Split(value, ",")
}

func getEnvCIDR(key string) ([]string, error) {
	cidr, err := netfilter.NormalizeCIDR(getEnvList(key, []string{}))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", key, err)
	}
	return cidr, nil
}

// getTrustedProxies loads TRUSTED_PROXIES and rejects catch-all prefixes
// (0.0.0.0/0, ::/0): trusting every address would bring back unconditional
// trust of X-Forwarded-For / X-Real-IP.
func getTrustedProxies() ([]string, error) {
	proxies, err := getEnvCIDR("TRUSTED_PROXIES")
	if err != nil {
		return nil, err
	}
	for _, proxy := range proxies {
		if netip.MustParsePrefix(proxy).Bits() == 0 {
			return nil, fmt.Errorf("TRUSTED_PROXIES: %q would trust every client, list the proxy addresses instead", proxy)
		}
	}
	return proxies, nil
}
