package repo

import "github.com/redis/go-redis/v9"

type IPStatsRepoI interface {
	UpSert()
}

type IPStatsRepo struct {
	redisClient *redis.Client
}

func NewIPStatRepo(redisClient *redis.Client) *IPStatsRepo {
	return &IPStatsRepo{
		redisClient: redisClient,
	}
}

func (repo *IPStatsRepo) UpSert() {
	return
}
