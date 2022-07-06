package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-datawrapper-adapter
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	DatawrapperUIURL           string        `envconfig:"DATAWRAPPER_UI_URL"`
	DatawrapperAPIURL          string        `envconfig:"DATAWRAPPER_API_URL"`
	DatawrapperAPIToken        string        `envconfig:"DATAWRAPPER_API_TOKEN"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                   ":28400",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		DatawrapperUIURL:           "https://app.datawrapper.de",
		DatawrapperAPIURL:          "https://api.datawrapper.de",
	}

	return cfg, envconfig.Process("", cfg)
}
