package backend

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
)

type BackendServer struct {
	URL          *url.URL
	Alive        bool
	ReverseProxy *httputil.ReverseProxy
	mux          sync.RWMutex
	activeConns  int64
}

func NewBackendServer(rawURL string) (*BackendServer, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(parsedURL)

	// Customize the reverse proxy to track connections
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Real-IP", getClientIP(req))
	}

	server := &BackendServer{
		URL:          parsedURL,
		Alive:        true,
		ReverseProxy: proxy,
		activeConns:  0,
	}

	// Wrap the reverse proxy to track active connections
	originalErrorHandler := proxy.ErrorHandler
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		atomic.AddInt64(&server.activeConns, -1)
		if originalErrorHandler != nil {
			originalErrorHandler(w, r, err)
		}
	}

	return server, nil
}

func (b *BackendServer) Serve(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&b.activeConns, 1)
	defer atomic.AddInt64(&b.activeConns, -1)

	b.ReverseProxy.ServeHTTP(w, r)
}

func (b *BackendServer) SetAlive(alive bool) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.Alive = alive
}

func (b *BackendServer) IsAlive() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.Alive
}

func (b *BackendServer) GetActiveConnections() int64 {
	return atomic.LoadInt64(&b.activeConns)
}
