package config

import (
	"time"

	"github.com/ONSdigital/dp-datawrapper-adapter/proxy"
	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-datawrapper-adapter
type Config struct {
	BindAddr                       string            `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout        time.Duration     `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval            time.Duration     `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout     time.Duration     `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	DatawrapperUIURL               proxy.URL         `envconfig:"DATAWRAPPER_UI_URL"`
	DatawrapperAPIURL              proxy.URL         `envconfig:"DATAWRAPPER_API_URL"`
	DatawrapperAPIToken            string            `envconfig:"DATAWRAPPER_API_TOKEN"`
	PermissionsAPIHost             string            `envconfig:"PERMISSIONS_API_HOST"`
	PermissionsCacheUpdateInterval time.Duration     `envconfig:"PERMISSIONS_CACHE_UPDATE_INTERVAL"`
	PermissionsMaxCacheTime        time.Duration     `envconfig:"PERMISSIONS_MAX_CACHE_TIME"`
	JWTVerificationPublicKeys      map[string]string `envconfig:"JWT_VERIFICATION_PUBLIC_KEYS"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                       ":28400",
		GracefulShutdownTimeout:        5 * time.Second,
		HealthCheckInterval:            30 * time.Second,
		HealthCheckCriticalTimeout:     90 * time.Second,
		PermissionsAPIHost:             "localhost:25400",
		PermissionsCacheUpdateInterval: time.Minute,
		PermissionsMaxCacheTime:        time.Minute * 5,
		JWTVerificationPublicKeys: map[string]string{
			"GHB723n83jw=": "MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA0TpTemKodQNChMNj1f/NF19nM",
			"HUJB8hw29js=": "MIICIjANBgkqBUHJHUJOIOIJIOH&*B(IHUGYCgKCAgEA0TpTemKodQNChMNj1f/NF19nM",
		},
	}

	return cfg, envconfig.Process("", cfg)
}
