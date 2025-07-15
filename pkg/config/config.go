package config

import (
	"fmt"
	"net/url"
	"os"
	"slices"
	"strconv"
)

const (
	// Env_dev is the development environment stack name.
	Env_dev = "dev"
	// Env_next is the next environment stack name.
	Env_next = "next"
	// Env_main is the main environment stack name.
	Env_main = "main"
)

var (
	// ErrEnvVarNotSet is a sentinel error that indicates that an environment variable is not set.
	ErrEnvVarNotSet = fmt.Errorf("required environment variable not set")
	// ErrEnvVarNotSetOrNil is a sentinel error that indicates that an environment variable is nil.
	ErrEnvVarNil = fmt.Errorf("required environment variable is nil")
	// ErrEnvVarInvalid is a sentinel error that indicates that an environment variable is invalid.
	ErrEnvVarInvalid = fmt.Errorf("environment variable is invalid")
	// ErrEnvVarInvalidValue is a sentinel error that indicates that an environment variable mapping has an invalid value.
	ErrEnvVarsMappingsInvalid = fmt.Errorf("environment variables mapping is invalid")
)

// A Config is a configuration struct that holds the configuration for an application.
// It is populated from environment variables and is used to configure the application.
// All fields are required and must be set in the environment.
type Config struct {
	// RuntimeEnv is the environment in which the application is running.
	// It must be one of `dev`, `next`, or `main`.
	RuntimeEnv string
	// AppVersion is the version of the application.
	// It should be a semantic version with no prefix, e.g. `1.0.0`.
	AppVersion string
	// AppName is the name of the application.
	// It should be a short name with no spaces or special characters, e.g. `shoppingcart`.
	// It is most likely the project / VCS reposiory name.
	AppName string
	// AppNamespace is the namespace of the application.
	// It describes the higher level project.
	// It should be a short name with no spaces or special characters, e.g. `shop`.
	AppNamespace string
	// OtelEndpointUrl is the URL of the OpenTelemetry Collector.
	// It should use `grpc` scheme whenever possible.
	OtelEndpointUrl url.URL
	// OtelExporterCompression is the compression to use when sending via OpenTelemetry.
	// It should be set to `gzip` whenever possible.
	OtelExporterCompression string
	// HttpServePort is the port on which the HTTP server will listen.
	// It must be a valid port number, e.g. `8080`.
	HttpServePort int
	// HttpReadTimeout is the maximum duration before timing out reads of the request.
	// It is passed to `http.Server.ReadTimeout`, multiplied by time.Second.
	HttpReadTimeout int
	// HttpWriteTimeout is the maximum duration before timing out writes of the response.
	// It is passed to `http.Server.WriteTimeout`, multiplied by time.Second.
	HttpWriteTimeout int
}

// envVarConf is a struct that holds the configuration for an environment variable.
type envVarConf struct {
	// Key is the name of the environment variable to look up.
	Key string
	// Dest is the destination where the value of the environment variable will be stored.
	Dest any
	// ValidSet is a set of valid values for the environment variable.
	ValidSet any
}

// getFromEnvVar returns the value of an environment variable by looking it up in the environment,
// and an error if it is not set or invalid.
// It returns the value of the environment variable and an error if it is not set or invalid.
func getFromEnvVar(cfg envVarConf) (any, error) {
	env, set := os.LookupEnv(cfg.Key)
	if !set {
		return nil, fmt.Errorf("env var %s: %w", cfg.Key, ErrEnvVarNotSet)
	}

	if env == "" {
		return nil, fmt.Errorf("env var %s: %w", cfg.Key, ErrEnvVarNil)
	}

	var res any

	switch cfg.Dest.(type) {
	case *string:
		res = env
	case *int:
		n, err := strconv.Atoi(env)
		if err != nil {
			return nil, fmt.Errorf("converting %s to int failed: %w", env, err)
		}

		res = n
	case *url.URL:
		u, err := url.Parse(env)
		if err != nil {
			return nil, fmt.Errorf("parsing %s to url failed: %w", env, err)
		}

		res = *u
	default:
		return nil, fmt.Errorf("%s: %w", cfg.Key, ErrEnvVarsMappingsInvalid)
	}

	if cfg.ValidSet != nil {
		switch set := cfg.ValidSet.(type) {
		case []string:
			if !slices.Contains(set, env) {
				return nil, fmt.Errorf("env var %s: %v not in set %v: %w", cfg.Key, env, set, ErrEnvVarInvalid)
			}
		case []int:
			found := false

			for _, c := range set {
				if res == c {
					found = true

					break
				}
			}

			if !found {
				return nil, fmt.Errorf("env var %s: %w", cfg.Key, ErrEnvVarInvalid)
			}
		default:
			return nil, fmt.Errorf("env var %s: %w", cfg.Key, ErrEnvVarsMappingsInvalid)
		}
	}

	return res, nil
}

// NewConfig returns a new Config struct populated with values from environment variables, merging them with default values,
// and an error if any of the required environment variables are not set or invalid.
func NewConfig() (*Config, error) {
	conf := &Config{}

	envVarsConfig := []envVarConf{
		{
			Key:  EnvVarKeyRuntimeEnv,
			Dest: &conf.RuntimeEnv,
			ValidSet: []string{
				Env_dev,
				Env_next,
				Env_main,
			},
		},
		{
			Key:  EnvVarKeyAppVersion,
			Dest: &conf.AppVersion,
		},
		{
			Key:  EnvVarKeyAppName,
			Dest: &conf.AppName,
		},
		{
			Key:  EnvVarKeyAppNamespace,
			Dest: &conf.AppNamespace,
		},
		{
			Key:  EnvVarKeyOtelEndpointURL,
			Dest: &conf.OtelEndpointUrl,
		},
		{
			Key:      EnvVarKeyOtelExporterCompression,
			Dest:     &conf.OtelExporterCompression,
			ValidSet: []string{"gzip"},
		},
		{
			Key:  EnvVarKeyHTTPServePort,
			Dest: &conf.HttpServePort,
		},
		{
			Key:  EnvVarKeyHTTPReadTimeout,
			Dest: &conf.HttpReadTimeout,
		},
		{
			Key:  EnvVarKeyHTTPWriteTimeout,
			Dest: &conf.HttpWriteTimeout,
		},
	}

	for _, env := range envVarsConfig {
		switch d := env.Dest.(type) {
		case *string:
			res, err := getFromEnvVar(env)
			if err != nil {
				return nil, err
			}

			s, ok := res.(string)
			if !ok {
				return nil, ErrEnvVarsMappingsInvalid
			}

			*d = s
		case *int:
			res, err := getFromEnvVar(env)
			if err != nil {
				return nil, err
			}

			i, ok := res.(int)
			if !ok {
				return nil, ErrEnvVarsMappingsInvalid
			}

			*d = i
		case *url.URL:
			res, err := getFromEnvVar(env)
			if err != nil {
				return nil, err
			}

			u, ok := res.(url.URL)
			if !ok {
				return nil, ErrEnvVarsMappingsInvalid
			}

			*d = u
		default:
			return nil, fmt.Errorf("%s: %w", env.Key, ErrEnvVarsMappingsInvalid)
		}
	}

	return conf, nil
}
