package repo

import (
	"log"
	"net/http"

	"github.com/mrdolev/sieve-go/utils"
)

type IPStatsHandlerI interface {
	UpSert(responseWriter http.ResponseWriter, request *http.Request)
}

type IPStatsHandler struct {
	service IPStatsServiceI
}

func NewHandler(service IPStatsServiceI) *IPStatsHandler {
	return &IPStatsHandler{
		service: service,
	}
}

func (handler *IPStatsHandler) UpSert(responseWriter http.ResponseWriter, request *http.Request) {
	ipFromAddress := utils.GetIPAddress(request)
	if ipFromAddress == "" {
		log.Printf("could not retrieve data ip \n")
		return
	}
	log.Println("ip address", ipFromAddress)
	handler.service.UpSert()
}
