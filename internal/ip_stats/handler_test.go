package ip_stats_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	ipstats "github.com/mrdolev/sieve-hll-go/internal/ip_stats"
	"github.com/mrdolev/sieve-hll-go/mocks"
	"github.com/mrdolev/sieve-hll-go/utils"
	"go.uber.org/mock/gomock"
)

const handlerIP string = "192.168.0.1"

func TestIPStatsHandlerCollect(t *testing.T) {
	tests := []struct {
		name       string
		headers    map[string]string
		wantIP     string
		statusCode int
		filter     bool
		invalid    bool
	}{
		{
			name:       "uses real ip header",
			headers:    map[string]string{utils.REAL_IP: handlerIP},
			wantIP:     handlerIP,
			statusCode: http.StatusOK,
			filter:     false,
		},
		{
			name: "prefers forwarded for header",
			headers: map[string]string{
				utils.FORWARD: "10.0.0.1, 10.0.0.2",
				utils.REAL_IP: handlerIP,
			},
			wantIP:     "10.0.0.1",
			statusCode: http.StatusOK,
			filter:     false,
		},
		{
			name:       "skip filtered ip",
			headers:    map[string]string{utils.REAL_IP: handlerIP},
			wantIP:     handlerIP,
			statusCode: http.StatusNoContent,
			filter:     true,
		},
		{
			name:       "rejects invalid ip",
			headers:    map[string]string{utils.REAL_IP: "not-an-ip"},
			statusCode: http.StatusBadRequest,
			invalid:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			ipStatsService := mocks.NewMockIPStatsServiceI(controller)
			netFilter := mocks.NewMockNetFilterI(controller)

			request := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			for key, value := range tt.headers {
				request.Header.Set(key, value)
			}

			if !tt.invalid {
				netFilter.EXPECT().Filter(tt.wantIP).Return(tt.filter).Times(1)
				ipStatsService.EXPECT().GetFilter().Return(netFilter).Times(1)
			}

			if !tt.invalid && !tt.filter {
				ipStatsService.EXPECT().Collect(request.Context(), tt.wantIP).Times(1)
			}

			handler := ipstats.NewHandler(ipStatsService)
			responseWriter := httptest.NewRecorder()

			handler.Collect(responseWriter, request)

			if responseWriter.Code != tt.statusCode {
				t.Fatalf("unexpected status code: got %d want %d", responseWriter.Code, tt.statusCode)
			}
		})
	}
}
