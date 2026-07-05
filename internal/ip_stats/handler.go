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
	ipFromAddress, err := utils.GetIPv4fromRequestRemoteAddr(request)
	if err != nil {
		log.Printf("%s", err.Error())
	}
	ipFromHeader, err := utils.GetIPv4FromRequestHeader(request)
	if err != nil {
		log.Printf("error to retrieve ip from request header")
	}
	if ipFromAddress != "" {
		responseWriter.Header().Add("From-Address", ipFromAddress)
		log.Println("fromAddress", ipFromAddress)
	}
	if ipFromHeader != "" {
		responseWriter.Header().Add("From-Header", ipFromHeader)
		log.Println("fromHeader", ipFromHeader)
	}
	handler.service.UpSert()
}
