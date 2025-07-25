package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Global struct {
	// Server configuration
	Server Server
	// Runtime configuration
	Runtime Runtime
	// Observability configuration
	Observability Observability
}

type Server struct {
	// Server bind address
	ListenAddr string `split_words:"true" default:"[::]:8080"`
	// HTTP read timeout
	ReadTimeout time.Duration `split_words:"true" default:"15s"`
	// HTTP write timeout
	WriteTimeout time.Duration `split_words:"true" default:"15s"`
	// HTTP idle timeout
	IdleTimeout time.Duration `split_words:"true" default:"60s"`
	// Proxy header for forwarded entity
	ProxyHeader string `split_words:"true" default:"X-Forwarded-For"`
}

type Runtime struct {
	// Environment the app is running in
	Environment string `split_words:"true" required:"true"`
	// Application version
	AppVersion string `split_words:"true" required:"true"`
	// Application name
	AppName string `split_words:"true" required:"true"`
	// Application namespace
	AppNamespace string `split_words:"true" required:"true"`
}

type Observability struct {
	// Address of OpenTelemetry endpoint where to send telemetry
	EndpointURL string `split_words:"true" required:"true"`
	// Compression to use when sending telemetry
	ExporterCompression string `split_words:"true"                 default:"gzip"`
	// Percentage of request to sample for tracing
	TracingSamplePercent int `split_words:"true"                 default:"100"`
}

// Load loads configuration from environment variables
func Load() (*Global, error) {
	var cfg Global
	err := envconfig.Process("kema", &cfg)
	if err != nil {
		return nil, fmt.Errorf("can't process config: %w", err)
	}
	return &cfg, nil
}

// IsLocalEnvironment returns whether the application in running in local-development environment
func IsLocalEnvironment(env string) bool {
	return env == "dev"
}
