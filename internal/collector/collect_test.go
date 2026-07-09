package collector

import (
	"context"
	"reflect"
	"testing"
	"time"
)

type fakeRepository struct {
	calls chan []string
}

func (repo *fakeRepository) PFAddBatch(ctx context.Context, ips []string) error {
	batch := append([]string(nil), ips...)

	select {
	case repo.calls <- batch:
	default:
	}

	return nil
}

func waitForBatch(t *testing.T, calls <-chan []string) []string {
	t.Helper()

	select {
	case batch := <-calls:
		return batch
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for batch flush")
		return nil
	}
}

func TestCollectorFlushesWhenBatchSizeReached(t *testing.T) {
	repo := &fakeRepository{calls: make(chan []string, 10)}
	collector := NewCollector(repo, 2, time.Hour)

	collector.Enqueue("192.168.0.1")
	collector.Enqueue("192.168.0.2")

	got := waitForBatch(t, repo.calls)
	want := []string{"192.168.0.1", "192.168.0.2"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected batch: got %v want %v", got, want)
	}

	collector.Close()
}

func TestCollectorResetsAfterFlush(t *testing.T) {
	repo := &fakeRepository{calls: make(chan []string, 10)}
	collector := NewCollector(repo, 2, time.Hour)
	defer collector.Close()

	collector.Enqueue("192.168.0.1")
	collector.Enqueue("192.168.0.2")
	first := waitForBatch(t, repo.calls)

	collector.Enqueue("192.168.0.3")
	collector.Enqueue("192.168.0.4")
	second := waitForBatch(t, repo.calls)

	if !reflect.DeepEqual(first, []string{"192.168.0.1", "192.168.0.2"}) {
		t.Fatalf("unexpected first batch: got %v", first)
	}

	if !reflect.DeepEqual(second, []string{"192.168.0.3", "192.168.0.4"}) {
		t.Fatalf("unexpected second batch: got %v", second)
	}
}

func TestCollectorFlushesOnClose(t *testing.T) {
	repo := &fakeRepository{calls: make(chan []string, 10)}
	collector := NewCollector(repo, 2, time.Hour)

	collector.Enqueue("192.168.0.1")
	collector.Close()

	got := waitForBatch(t, repo.calls)
	want := []string{"192.168.0.1"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected batch on close: got %v want %v", got, want)
	}
}
