package collector

import (
	"context"
	"log"
	"time"

	"github.com/mrdolev/sieve-go/internal/storage"
)

type CollectorI interface {
	Enqueue(ip string)
	Close()
}

type Collector struct {
	repo      storage.Repository
	bufferCh  chan string
	batchSize int
	interval  time.Duration
	doneCh    chan struct{}
}

func NewCollector(repo storage.Repository, batchSize int, interval time.Duration) *Collector {
	return &Collector{
		repo:      repo,
		bufferCh:  make(chan string, 10000),
		batchSize: batchSize,
		interval:  interval,
		doneCh:    make(chan struct{}),
	}
}

func (collector *Collector) Enqueue(ip string) {
	collector.bufferCh <- ip
}

func (collector *Collector) Close() {
	close(collector.doneCh)
}

func (collector *Collector) run() {
	ticker := time.NewTicker(collector.interval)
	defer ticker.Stop()

	batch := make([]string, 0, collector.batchSize)
	for {
		select {
		case ip := <-collector.bufferCh:
			batch = append(batch, ip)
			if len(batch) >= collector.batchSize {
				collector.flush(batch)
			}
		case <-ticker.C:
			if len(batch) > 0 {
				collector.flush(batch)
				// reset
				batch = batch[:0]
			}
		case <-collector.doneCh:
			if len(batch) > 0 {
				collector.flush(batch)
			}
			return
		}
	}
}

func (collector *Collector) flush(batch []string) {
	if err := collector.repo.PFAddBatch(context.Background(), batch); err != nil {
		log.Printf("error to store information %s", err.Error())
		return
	}
	log.Printf("the information are stored correctly %d", len(batch))
	return
}
