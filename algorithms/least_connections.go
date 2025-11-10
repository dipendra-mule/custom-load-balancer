package algorithms

import (
	"net/http"

	"github.com/dipendra-mule/custom-load-balancer/backend"
)

type LeastConnections struct{}

func (lc *LeastConnections) SelectBackend(pool *backend.BackendPool, r *http.Request) *backend.BackendServer {
	servers := pool.GetHealthyServers()
	if len(servers) == 0 {
		return nil
	}

	var selected *backend.BackendServer
	minConnections := int64(^uint64(0) >> 1) // max int64

	for _, server := range servers {
		conns := server.GetActiveConnections()
		if conns < minConnections {
			minConnections = conns
			selected = server
		}
	}

	return selected
}

func (lc *LeastConnections) Name() string {
	return "least-connections"
}
