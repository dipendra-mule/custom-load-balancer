package algorithms

import (
	"hash/fnv"
	"net/http"
	"strings"

	"github.com/dipendra-mule/custom-load-balancer/backend"
)

type IPHash struct{}

func (ih *IPHash) SelectBackend(pool *backend.BackendPool, r *http.Request) *backend.BackendServer {
	server := pool.GetHealthyServers()

	if len(server) == 0 {
		return nil
	}

	// get client IP
	ip := getClientIP(r)

	// create hash
	h := fnv.New32a()
	h.Write([]byte(ip))
	hash := h.Sum32()

	index := hash % uint32(len(server))
	return server[index]
}

func getClientIP(r *http.Request) string {
	// check x-forwarded-for header
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// check x-real-ip
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	// fall back to remote address
	ip := strings.Split(r.RemoteAddr, ":")[0]
	return ip
}

func (ih *IPHash) Name() string {
	return "ip-hash"
}
