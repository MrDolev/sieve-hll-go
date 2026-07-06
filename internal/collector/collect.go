package collector

import (
	"time"

	ip_stats "github.com/mrdolev/sieve-go/internal/ip_stats"
)

type CollectorI interface {
	Enqueue(ip string)
	Close()
}

type Collector struct {
	repo      ip_stats.IPStatsRepo
	ch        chan string
	batchSize int
	interval  time.Duration
	done      chan struct{}
}

func NewCollector() *Collector {
	return &Collector{}
}

func (collector *Collector) Enqueue(ip string) {
	return
}

func (collector *Collector) Close() {
	return
}

func (collector *Collector) run() {
	return
}

func (collector *Collector) flush() {
	return
}
