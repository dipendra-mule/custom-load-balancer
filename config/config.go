package config

import "time"

type LoadBalancerConfig struct {
	Port           int           `yaml:"port"`
	Algorithm      string        `yaml:"algorithm"`
	HealthCheck    HealthCheck   `yaml:"health_check"`
	TLS            TLSConfig     `yaml:"tls"`
	StickySessions StickyConfig  `yaml:"sticky_sessions"`
	BackendServers []Backend     `yaml:"backend_servers"`
	Observability  Observability `yaml:"observability"`
}

type HealthCheck struct {
	Enabled      bool          `yaml:"enabled"`
	Interval     time.Duration `yaml:"interval"`
	Timeout      time.Duration `yaml:"timeout"`
	Path         string        `yaml:"path"`
	SuccessCodes []int         `yaml:"success_codes"`
}

type TLSConfig struct {
	Enabled      bool          `yaml:"enabled"`
	Interval     time.Duration `yaml:"interval"`
	Timeout      time.Duration `yaml:"timeout"`
	Path         string        `yaml:"path"`
	SuccessCodes []int         `yaml:"success_codes"`
}

type StickyConfig struct {
	Enabled    bool          `yaml:"enabled"`
	Duration   time.Duration `yaml:"duration"`
	CookieName string        `yaml:"cookie_name"`
}

type Backend struct {
	URL    string `yaml:"url"`
	Weight int    `yaml:"weight"`
}

type Observability struct {
	MetricsEnabled bool   `yaml:"metrics_enabled"`
	LogLevel       string `yaml:"log_level"`
	TraceEnabled   bool   `yaml:"trace_enabled"`
}
