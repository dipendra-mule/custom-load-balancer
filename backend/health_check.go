package backend

import (
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
