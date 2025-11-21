package backend

import (
	"net/http"
	"net/url"
	"time"
)

type HealthChecker struct {
	pool     *BackendPool
	config   HealthCheckConfig
	stopChan chan bool
}

type HealthCheckConfig struct {
	Enabled      bool
	Interval     time.Duration
	Timeout      time.Duration
	Path         string
	SuccessCodes []int
}

func NewHealthChecker(pool *BackendPool, config HealthCheckConfig) *HealthChecker {
	return &HealthChecker{
		pool:     pool,
		config:   config,
		stopChan: make(chan bool),
	}
}

func (hc *HealthChecker) Start() {
	ticker := time.NewTicker(hc.config.Interval)

	go func() {
		for {
			select {
			case <-ticker.C:
				hc.checkHealth()
			case <-hc.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

func (hc *HealthChecker) Stop() {
	close(hc.stopChan)
}

func (hc *HealthChecker) checkHealth() {
	servers := hc.pool.GetServers()

	for _, server := range servers {
		go func(s *BackendServer) {
			alive := hc.isBackendAlive(s.URL)
			hc.pool.MarkBackendStatus(s.URL, alive)
		}(server)
	}
}

func (hc *HealthChecker) isBackendAlive(backendURL *url.URL) bool {
	client := http.Client{
		Timeout: hc.config.Timeout,
	}

	healthURL := *backendURL
	healthURL.Path = hc.config.Path

	resp, err := client.Get(healthURL.String())
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// Check if status code is in success codes
	for _, code := range hc.config.SuccessCodes {
		if resp.StatusCode == code {
			return true
		}
	}

	return false
}
