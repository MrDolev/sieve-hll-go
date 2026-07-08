package ip_stats

import (
	"context"
	"log"

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
	log.Printf("start to collect from service - %s \n", ip)
	service.collector.Enqueue(ip)
}
