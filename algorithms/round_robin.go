package algorithms

import (
	"net/http"
	"sync/atomic"

	"github.com/dipendra-mule/custom-load-balancer/backend"
)

type RoundRobin struct {
	counter uint64
}

func (rr *RoundRobin) SelectBackend(pool *backend.BackendPool, r *http.Request) *backend.BackendServer {
	servers := pool.GetHealthyServers()
	if len(servers) == 0 {
		return nil
	}

	index := atomic.AddUint64(&rr.counter, 1) % uint64(len(servers))
	return servers[index]
}

func (rr *RoundRobin) Name() string {
	return "round-robin"
}
