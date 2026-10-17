package ip_stats

import (
	"log"
	"net/http"
	"net/netip"

	"github.com/mrdolev/sieve-hll-go/utils"
)

type IPStatsHandlerI interface {
	Collect(responseWriter http.ResponseWriter, request *http.Request)
}

type IPStatsHandler struct {
	service IPStatsServiceI
}

func NewHandler(service IPStatsServiceI) *IPStatsHandler {
	return &IPStatsHandler{
		service: service,
	}
}

func (handler *IPStatsHandler) Collect(responseWriter http.ResponseWriter, request *http.Request) {
	ipFromAddress := utils.GetIPAddress(request)
	if ipFromAddress == "" {
		log.Printf("could not retrieve data ip \n")
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := netip.ParseAddr(ipFromAddress); err != nil {
		log.Printf("invalid ip address %q \n", ipFromAddress)
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	if handler.service.GetFilter().Filter(ipFromAddress) {
		log.Printf("ip address %s is filtered", ipFromAddress)
		responseWriter.WriteHeader(http.StatusNoContent)
		return
	}
	log.Println("ip address", ipFromAddress)
	handler.service.Collect(request.Context(), ipFromAddress)
}
