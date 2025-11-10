package algorithms

import (
	"net/http"
	"sync/atomic"
)

type RoundRobin struct {
	counter int
}

func (rr *RoundRobin) SelectBackend(pool *backend.BackendPool, r *http.Request) *backend.BackendServer {
	servers := pool.GetHealthyServers()
	if len(servers) == 0 {
		return nil
	}

	index := atomic.AddUint64(&rr.counter, 1) % uint64(len(servers))
	return servers[index]
}

func (r *RoundRobin) Name() string {
	return "round-robin"
}
