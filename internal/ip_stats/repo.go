package ip_stats

import (
	"context"
	"time"

	"github.com/mrdolev/sieve-hll-go/internal/storage"
)

type IPStatsRepo struct {
	storageClient *storage.RedisClient
}

func NewIPStatsRepo(storageClient *storage.RedisClient) *IPStatsRepo {
	return &IPStatsRepo{
		storageClient: storageClient,
	}
}

func (ipStatsRepo *IPStatsRepo) PFAddBatch(ctx context.Context, ips []string) error {
	args := make([]interface{}, len(ips))

	for index, value := range ips {
		args[index] = value
	}

	key := "client requests" + time.Now().Format("2006-01-02")

	return ipStatsRepo.storageClient.Client().PFAdd(ctx, key, args...).Err()

}
