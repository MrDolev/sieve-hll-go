package ip_stats

import (
	"context"
	"log"

	"github.com/mrdolev/sieve-hll-go/internal/collector"
	netfilter "github.com/mrdolev/sieve-hll-go/internal/net_filter"
)

type IPStatsServiceI interface {
	Collect(ctx context.Context, ip string)
	GetFilter() netfilter.NetFilterI
}

type IPStatsService struct {
	collector collector.CollectorI
	filter    netfilter.NetFilterI
}

func NewIPStatsService(collector collector.CollectorI, filter netfilter.NetFilterI) *IPStatsService {
	return &IPStatsService{
		collector: collector,
		filter:    filter,
	}
}

func (service *IPStatsService) Collect(ctx context.Context, ip string) {
	log.Printf("start to collect from service - %s \n", ip)
	service.collector.Enqueue(ip)
}

func (service *IPStatsService) GetFilter() netfilter.NetFilterI {
	return service.filter
}
