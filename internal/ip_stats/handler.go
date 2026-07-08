package ip_stats

import (
	"log"
	"net/http"

	"github.com/mrdolev/sieve-go/utils"
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
		return
	}
	log.Println("ip address", ipFromAddress)
	handler.service.Collect(request.Context(), ipFromAddress)
}
