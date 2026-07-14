package ip_stats

import (
	"context"
	"log"

	"github.com/mrdolev/sieve-hll-go/internal/collector"
)

type IPStatsServiceI interface {
	Collect(ctx context.Context, ip string)
}

type IPStatsService struct {
	collector collector.CollectorI
}

func NewIPStatsService(collector collector.CollectorI) *IPStatsService {
	return &IPStatsService{
		collector: collector,
	}
}

func (service *IPStatsService) Collect(ctx context.Context, ip string) {
	log.Printf("start to collect from service - %s \n", ip)
	service.collector.Enqueue(ip)
}
