// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/kemadev/go-framework/pkg/semver"
)

const ConfigurationEnvVarPrefix = "kema"

var (
	ErrCantRedactDBConnURL = errors.New("database connection URL is not redactable")
	ErrCantMarshallConfig  = errors.New("config is not marshallable to JSON")
)

var (
	ErrVariableRequired  = errors.New("environment variable required")
	ErrVariableMalformed = errors.New("environment malformed")
)

// Global is the server configuration struct.
// Values are populated from environment variables nammed after
// their relative position in the struct with [ConfigurationEnvVarPrefix] as prefix, using SCREAMING_SNAKE_CASE.
// e.g. [Global.Observability.EndpointURL] is populated from environment variable `[ConfigurationEnvVarPrefix]_OBSERVABILITY_ENDPOINT_URL`.
type Global struct {
	// Server holds the HTTP server configuration
	Server Server `required:"true"`
	// Runtime holds the runtime configuration
	Runtime Runtime `required:"true"`
	// Observability holds the observability configuration
	Observability Observability `required:"true"`
	// Client holds the clients configurations
	Client Client `required:"false"`
}

// Server holds the HTTP server configuration.
type Server struct {
	// BindAddr is the server bind addressfor the HTTP server
	BindAddr string `default:"[::]"      required:"true"`
	// BindPort is the server bind portfor the HTTP server
	BindPort int `default:"8080"      required:"true"`
	// ReadTimeout is the HTTP read timeout for the HTTP server
	ReadTimeout time.Duration `default:"15s"       required:"true"`
	// WriteTimeout is the HTTP write timeout for the HTTP server
	WriteTimeout time.Duration `default:"15s"       required:"true"`
	// IdleTimeout is the HTTP idle timeout for the HTTP server
	IdleTimeout time.Duration `default:"60s"       required:"true"`
	// ProxyHeader is the proxy header for forwarded entity
	ProxyHeader string `default:"Forwarded" required:"true"`
	// ShutdownGracePeriod is the grace period to give the server before canceling contexits t upon shutdown
	ShutdownGracePeriod time.Duration `default:"5s"        required:"true"`
}

// Runtime holds the runtime configuration.
type Runtime struct {
	// Environment the app is running in
	Environment string `required:"true"`
	// Application version
	AppVersion semver.Version `required:"true"`
	// Application name
	AppName string `required:"true"`
	// Application namespace
	AppNamespace string `required:"true"`
}

// Observability holds the observability configuration.
type Observability struct {
	// Address of OpenTelemetry endpoint where to send telemetry
	EndpointURL url.URL `required:"true"`
	// Compression to use when sending telemetry
	ExporterCompression string `required:"true" default:"gzip"`
	// Percentage of request to sample for tracing
	TracingSamplePercent int `required:"true" default:"100"`
	// Interval between metrics exports, in seconds
	MetricsExportInterval time.Duration `required:"true" default:"15s"`
	// ShutdownGracePeriod is the grace period to give the instrumentation before canceling its context upon shutdown
	ShutdownGracePeriod time.Duration `required:"true" default:"5s"`
}

// Client holds the clients configurations.
type Client struct {
	// Database holds database configuration
	Database DatabaseConfig `required:"false"`
	// Cache holds cache configuration
	Cache CacheConfig `required:"false"`
	// ObjectStorage holds object storage configuration
	ObjectStorage ObjectStorageConfig `required:"false"`
}

type DatabaseConfig struct {
	// Connection URL used to connect to the database
	ConnectionURL url.URL `required:"false"`
}

type CacheConfig struct {
	ClientAddress         []string      `required:"true"`
	ShardsRefreshInterval time.Duration `required:"true"  default:"120s"`
	SentinelMasterSet     string        `required:"false"`
	Username              string        `required:"true"`
	Password              string        `required:"true"`
}

type ObjectStorageConfig struct {
	EndpointAddressAddress string `required:"true"`
	AccessKeyID            string `required:"true"`
	SecretAccessKey        string `required:"true"`
	SSL                    bool   `required:"true"`
}

// Manager handles configuration loading and caching.
type Manager struct {
	once   sync.Once
	config *Global
	err    error
}

// NewManager creates a new configuration manager.
func NewManager() *Manager {
	return &Manager{}
}

// Load loads configuration from environment variables
// On first call, it loads the configuration from environment variables.
// Subsequent calls return the cached configuration.
func (m *Manager) Load() (*Global, error) {
	m.once.Do(func() {
		var conf Global

		err := load(ConfigurationEnvVarPrefix, &conf)
		if err != nil {
			m.err = fmt.Errorf("can't process config: %w", err)

			return
		}

		m.config = &conf
	})

	if m.err != nil {
		return nil, m.err
	}

	return m.config, nil
}

// Reset clears the cached configuration and allows Load() to reload it.
// This is primarily useful for testing scenarios.
func (m *Manager) Reset() {
	m.once = sync.Once{}
	m.config = nil
	m.err = nil
}

// Get returns the loaded configuration or loads it if not already loaded.
func (m *Manager) Get() (*Global, error) {
	return m.Load()
}

// EnvLocalValue is the value of the environment variable backing [Runtime.Environment] which is used to denote a local development environment.
const EnvLocalValue = "dev"

// IsLocalEnvironment returns whether the application in running in local-development environment.
func (conf *Runtime) IsLocalEnvironment() bool {
	return conf.Environment == EnvLocalValue
}

// load processes configuration from environment variables with the given prefix.
func load(prefix string, cfg any) error {
	return processStruct(prefix, reflect.ValueOf(cfg).Elem(), "", true)
}

// processStruct recursively processes struct fields.
func processStruct(prefix string, v reflect.Value, parentPath string, parentRequired bool) error {
	t := v.Type()

	for i := range v.NumField() {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		fieldName := fieldType.Name
		envVarName := buildEnvVarName(prefix, parentPath, fieldName)

		switch field.Kind() {
		case reflect.Struct:
			// Check if this struct field is required
			structRequired := fieldType.Tag.Get("required") == "true"

			// If parent is not required, this struct is also not required
			if !parentRequired {
				structRequired = false
			}

			err := processStruct(prefix, field, buildPath(parentPath, fieldName), structRequired)
			if err != nil {
				return fmt.Errorf("error processing struct field %s: %w", fieldName, err)
			}
		default:
			err := processField(field, fieldType, envVarName, parentRequired)
			if err != nil {
				return fmt.Errorf("error processing field %s: %w", fieldType.Name, err)
			}
		}
	}

	return nil
}

// processField processes a single struct field.
func processField(
	field reflect.Value,
	fieldType reflect.StructField,
	envVarName string,
	parentRequired bool,
) error {
	defaultValue := fieldType.Tag.Get("default")
	fieldRequired := fieldType.Tag.Get("required") == "true"

	// If parent struct is not required, this field is also not required
	required := fieldRequired && parentRequired

	envValue := os.Getenv(envVarName)

	if envValue == "" && defaultValue != "" {
		envValue = defaultValue
	}

	if required && envValue == "" {
		return fmt.Errorf("%s: %w", envVarName, ErrVariableRequired)
	}

	if envValue == "" {
		return nil
	}

	return setFieldValue(field, envValue, envVarName)
}

// setFieldValue sets the field value based on its type.
func setFieldValue(field reflect.Value, value string, envVarName string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Slice:
		if field.Type().Elem().Kind() == reflect.String {
			return setStringSlice(field, value)
		}
		return fmt.Errorf(
			"%s - unsupported slice type %s: %w",
			envVarName,
			field.Type(),
			ErrVariableMalformed,
		)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Type() == reflect.TypeOf(time.Duration(0)) {
			duration, err := time.ParseDuration(value)
			if err != nil {
				return fmt.Errorf("%s - %s: %w", envVarName, value, ErrVariableMalformed)
			}

			field.SetInt(int64(duration))
		} else {
			intVal, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("%s - %s: %w", envVarName, value, ErrVariableMalformed)
			}

			field.SetInt(intVal)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("%s - %s: %w", envVarName, value, ErrVariableMalformed)
		}

		field.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("%s - %s: %w", envVarName, value, ErrVariableMalformed)
		}

		field.SetFloat(floatVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("%s - %s: %w", envVarName, value, ErrVariableMalformed)
		}

		field.SetBool(boolVal)
	case reflect.Struct:
		switch field.Type() {
		case reflect.TypeOf(url.URL{}):
			parsedURL, err := url.Parse(value)
			if err != nil {
				return fmt.Errorf(
					"%s - invalid URL %s: %w",
					envVarName,
					value,
					ErrVariableMalformed,
				)
			}
			field.Set(reflect.ValueOf(*parsedURL))
		case reflect.TypeOf(semver.Version{}):
			parsedVersion, err := semver.Parse(value)
			if err != nil {
				return fmt.Errorf(
					"%s - invalid version %s: %w",
					envVarName,
					value,
					ErrVariableMalformed,
				)
			}
			field.Set(reflect.ValueOf(parsedVersion))
		default:
			return fmt.Errorf(
				"%s - unsupported struct type %s: %w",
				envVarName,
				field.Type(),
				ErrVariableMalformed,
			)
		}
	default:
		return fmt.Errorf("%s - %s: %w", field.Kind(), envVarName, ErrVariableMalformed)
	}
	return nil
}

func setStringSlice(field reflect.Value, value string) error {
	if value == "" {
		field.Set(reflect.MakeSlice(field.Type(), 0, 0))
		return nil
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	slice := reflect.MakeSlice(field.Type(), len(result), len(result))
	for i, str := range result {
		slice.Index(i).SetString(str)
	}
	field.Set(slice)

	return nil
}

// buildEnvVarName builds the environment variable name from prefix and field path.
func buildEnvVarName(prefix, parentPath, fieldName string) string {
	parts := []string{CamelToScreamingSnake(prefix)}

	if parentPath != "" {
		parts = append(parts, CamelToScreamingSnake(parentPath))
	}

	parts = append(parts, CamelToScreamingSnake(fieldName))

	return strings.Join(parts, "_")
}

// buildPath builds the path for nested structs.
func buildPath(parentPath, fieldName string) string {
	if parentPath == "" {
		return fieldName
	}

	return parentPath + fieldName
}

// CamelToScreamingSnake converts camelCase to SCREAMING_SNAKE_CASE.
func CamelToScreamingSnake(str string) string {
	var result strings.Builder

	for pos, char := range str {
		if pos > 0 && unicode.IsUpper(char) {
			// Check if previous character was lowercase or if next character is lowercase
			prevChar := rune(str[pos-1])
			nextIsLower := pos+1 < len(str) && unicode.IsLower(rune(str[pos+1]))

			if unicode.IsLower(prevChar) || nextIsLower {
				result.WriteRune('_')
			}
		}

		result.WriteRune(unicode.ToLower(char))
	}

	return strings.ToUpper(result.String())
}

// SlogLevel return the appropriate [slog.Level] for given [Runtime].
func (conf *Runtime) SlogLevel() slog.Level {
	if conf.IsLocalEnvironment() {
		return slog.LevelDebug
	}

	return slog.LevelInfo
}

// Redact returns a modified version of conf, redacting sensible values
func (conf Global) Redact() (Global, error) {
	DBPasswd, present := conf.Client.Database.ConnectionURL.User.Password()
	if present {
		redactedPasswd := DBPasswd[0:(len(DBPasswd)/4)] + "..."
		redactedURLParts := strings.Split(conf.Client.Database.ConnectionURL.String(), DBPasswd)
		if len(redactedURLParts) != 2 {
			return Global{}, ErrCantRedactDBConnURL
		} else {
			redactedURLStr := strings.Join(redactedURLParts, redactedPasswd)
			redactedURL, err := url.Parse(redactedURLStr)
			if err != nil {
				return Global{}, ErrCantRedactDBConnURL
			}
			conf.Client.Database.ConnectionURL = *redactedURL
		}
	}
	cachePasswdLen := len(conf.Client.Cache.Password)
	if cachePasswdLen > 0 {
		conf.Client.Cache.Password = conf.Client.Cache.Password[0:(cachePasswdLen/4)] + "..."
	}

	ObjectStorageAccessKeyLen := len(conf.Client.ObjectStorage.SecretAccessKey)
	if ObjectStorageAccessKeyLen > 0 {
		conf.Client.ObjectStorage.SecretAccessKey = conf.Client.ObjectStorage.SecretAccessKey[0:(ObjectStorageAccessKeyLen/4)] + "..."
	}

	return conf, nil
}

func (conf Global) String() (string, error) {
	b, err := json.Marshal(conf)
	if err != nil {
		return "", ErrCantMarshallConfig
	}

	return string(b), nil
}
