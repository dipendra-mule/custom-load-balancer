package backend

import (
	"net/url"
	"sync"
)

type BackendPool struct {
	servers []*BackendServer
	mux     sync.RWMutex
}

func NewBackendPool() *BackendPool {
	return &BackendPool{
		servers: make([]*BackendServer, 0),
	}
}

func (bp *BackendPool) AddBackend(server *BackendServer) {
	bp.mux.Lock()
	defer bp.mux.Unlock()
	bp.servers = append(bp.servers, server)
}

func (bp *BackendPool) GetServers() []*BackendServer {
	bp.mux.RLock()
	defer bp.mux.RUnlock()
	return bp.servers
}

func (bp *BackendPool) GetHealthyServers() []*BackendServer {
	bp.mux.RLock()
	defer bp.mux.RUnlock()

	healthy := make([]*BackendServer, 0)
	for _, server := range bp.servers {
		if server.IsAlive() {
			healthy = append(healthy, server)
		}
	}
	return healthy
}

func (bp *BackendPool) MarkBackendStatus(serverURL *url.URL, alive bool) {
	bp.mux.Lock()
	defer bp.mux.Unlock()

	for _, server := range bp.servers {
		if server.URL.String() == serverURL.String() {
			server.SetAlive(alive)
			break
		}
	}
}
