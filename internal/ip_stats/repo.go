package ip_stats

import (
	"context"
	"time"

	"github.com/mrdolev/sieve-go/internal/storage"
)

type IPStatsRepo struct {
	storageClient *storage.RedisClient
}

func NewIPStatRepo(storageClient *storage.RedisClient) *IPStatsRepo {
	return &IPStatsRepo{
		storageClient: storageClient,
	}
}

func (repo *IPStatsRepo) PFAddBatch(ctx context.Context, ips []string) error {
	args := make([]interface{}, len(ips))

	for index, value := range ips {
		args[index] = value
	}

	key := "client requests" + time.Now().Format("2006-01-02")

	return repo.storageClient.Client().PFAdd(ctx, key, args...).Err()

}
