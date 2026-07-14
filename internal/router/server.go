package server

import (
	"net/http"

	repoIPStats "github.com/mrdolev/sieve-hll-go/internal/ip_stats"
)

type Router struct {
	handler   repoIPStats.IPStatsHandlerI
	serverMux *http.ServeMux
}

func NewRouter(handler repoIPStats.IPStatsHandlerI, server *http.ServeMux) *Router {
	return &Router{
		handler:   handler,
		serverMux: server,
	}
}

func (r *Router) Collect(path string) {
	r.serverMux.HandleFunc(path, r.handler.Collect)
}
