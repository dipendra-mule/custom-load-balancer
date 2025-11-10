package algorithms

import (
	"net/http"

	"github.com/dipendra-mule/custom-load-balancer/backend"
)

type LoadBalancerAlgorithm interface {
	SelectBackend(pool *backend.BackendPool, r *http.Request) *backend.BackendServers
	Name() string
}

type AlgorithmFactory struct {
	algorithms map[string]func() LoadBalancingAlgorithm
}

func NewAlgorithmFactory() *AlgorithmFactory {
	af := &AlgorithmFactory{
		algorithms: make(map[string]func() LoadBalancingAlgorithm),
	}

	af.Register("round-robin", func() LoadBalancerAlgorithm { return &RoundRobin{} })
	af.Register("least-connections", func() LoadBalancerAlgorithm { return &LeastConnections{} })
	af.Register("ip-hash", func() LoadBalancerAlgorithm { return &IPHash{} })

	return af
}

func (af *AlgorithmFactory) Register(name string, factory func() LoadBalancingAlgorithm) {
	af.algorithms[name] = factory
}

func (af *AlgorithmFactory) Get(name string) LoadBalancerAlgorithm {
	if factory, ok := af.algorithms[name]; ok {
		return factory()
	}
	return af.algorithms["round-robin"]() // default
}
