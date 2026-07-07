package repo

import (
	"context"

	"github.com/mrdolev/sieve-go/internal/collector"
	"github.com/mrdolev/sieve-go/internal/storage"
)

type IPStatsServiceI interface {
	Collect(ctx context.Context, ip string)
}

type IPStatsService struct {
	repo      storage.Repository
	collector collector.CollectorI
}

func NewIPStatsService(repo storage.Repository, collector collector.CollectorI) *IPStatsService {
	return &IPStatsService{
		repo:      repo,
		collector: collector,
	}
}

func (service *IPStatsService) Collect(ctx context.Context, ip string) {
	service.collector.Enqueue(ip)
}
