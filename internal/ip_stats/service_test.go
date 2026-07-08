package ip_stats_test

import (
	"context"
	"testing"

	ipstats "github.com/mrdolev/sieve-go/internal/ip_stats"
	"github.com/mrdolev/sieve-go/mocks"
	"go.uber.org/mock/gomock"
)

const serviceIP string = "192.168.0.1"

func TestIPStatsServiceCollect(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	collectorMock := mocks.NewMockCollectorI(controller)
	collectorMock.EXPECT().Enqueue(serviceIP).Times(1)

	service := ipstats.NewIPStatsService(nil, collectorMock)
	if service == nil {
		t.Fatal("expected service instance")
	}

	service.Collect(context.Background(), serviceIP)
}
