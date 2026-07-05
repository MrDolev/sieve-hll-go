package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/mrdolev/sieve-go/utils"
)

type ServerMux struct {
	port int
	path string
	mux  *http.ServeMux
}

func NewServerMux(port int, path string) *ServerMux {
	mux := http.NewServeMux()
	return &ServerMux{
		port: port,
		path: path,
		mux:  mux,
	}
}

func exploreIP(responseWriter http.ResponseWriter, request *http.Request) {
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
}

func (s *ServerMux) Serve() {
	s.mux.HandleFunc(s.path, exploreIP)
	err := http.ListenAndServe(fmt.Sprintf(":%s", strconv.Itoa(s.port)), s.mux)
	if err != nil {
		log.Fatalf("error to start server %s ", err.Error())
	}
}
