package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const ConfigurationEnvVarPrefix = "kema"

// Global is the server configuration struct.
// Values are populated from environment variables nammed after
// their relative position in the struct with [ConfigurationEnvVarPrefix] as prefix, using SCREAMING_SNAKE_CASE.
// e.g. [Global.Observability.EndpointURL] is populated from environment variable `[ConfigurationEnvVarPrefix]_OBSERVABILITY_ENDPOINT_URL`
type Global struct {
	// Server holds the HTTP server configuration
	Server Server
	// Runtime holds the runtime configuration
	Runtime Runtime
	// Observability holds the observability configuration
	Observability Observability
}

// Server holds the HTTP server configuration
type Server struct {
	// BindAddr is the server bind addressfor the HTTP server
	BindAddr string `required:"true" default:"[::]"`
	// BindPort is the server bind portfor the HTTP server
	BindPort int `required:"true" default:"8080"`
	// ReadTimeout is the HTTP read timeout for the HTTP server
	ReadTimeout time.Duration `required:"true" default:"15s"`
	// WriteTimeout is the HTTP write timeout for the HTTP server
	WriteTimeout time.Duration `required:"true" default:"15s"`
	// IdleTimeout is the HTTP idle timeout for the HTTP server
	IdleTimeout time.Duration `required:"true" default:"60s"`
	// ProxyHeader is the proxy header for forwarded entity
	ProxyHeader string `required:"true" default:"X-Forwarded-For"`
	// ShutdownGracePeriod is the grace period to give the server before canceling contexits t upon shutdown
	ShutdownGracePeriod time.Duration `required:"true" default:"5s"`
}

// Runtime holds the runtime configuration
type Runtime struct {
	// Environment the app is running in
	Environment string `required:"true"`
	// Application version
	AppVersion string `required:"true"`
	// Application name
	AppName string `required:"true"`
	// Application namespace
	AppNamespace string `required:"true"`
}

// Observability holds the observability configuration
type Observability struct {
	// Address of OpenTelemetry endpoint where to send telemetry
	EndpointURL string `required:"true"`
	// Compression to use when sending telemetry
	ExporterCompression string `required:"true" default:"gzip"`
	// Percentage of request to sample for tracing
	TracingSamplePercent int `required:"true" default:"100"`
	// Interval between metrics exports, in seconds
	MetricsExportInterval time.Duration `required:"true" default:"15s"`
	// ShutdownGracePeriod is the grace period to give the instrumentation before canceling its context upon shutdown
	ShutdownGracePeriod time.Duration `required:"true" default:"5s"`
}

// Load loads configuration from environment variables
func Load() (*Global, error) {
	var cfg Global
	err := load(ConfigurationEnvVarPrefix, &cfg)
	if err != nil {
		return nil, fmt.Errorf("can't process config: %w", err)
	}

	return &cfg, nil
}

// IsLocalEnvironment returns whether the application in running in local-development environment
func (cfg Runtime) IsLocalEnvironment() bool {
	return cfg.Environment == "dev"
}

// load processes configuration from environment variables with the given prefix
func load(prefix string, cfg interface{}) error {
	return processStruct(prefix, reflect.ValueOf(cfg).Elem(), "")
}

// processStruct recursively processes struct fields
func processStruct(prefix string, v reflect.Value, parentPath string) error {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
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
			err := processStruct(prefix, field, buildPath(parentPath, fieldName))
			if err != nil {
				return fmt.Errorf("failure processing struct field %s: %w", fieldName, err)
			}
		default:
			err := processField(field, fieldType, envVarName)
			if err != nil {
				return fmt.Errorf("failure processing field %s: %w", fieldType.Name, err)
			}
		}
	}

	return nil
}

// processField processes a single struct field
func processField(field reflect.Value, fieldType reflect.StructField, envVarName string) error {
	defaultValue := fieldType.Tag.Get("default")
	required := fieldType.Tag.Get("required") == "true"

	envValue := os.Getenv(envVarName)

	if envValue == "" && defaultValue != "" {
		envValue = defaultValue
	}

	if required && envValue == "" {
		return fmt.Errorf("required environment variable %s is not set", envVarName)
	}

	if envValue == "" {
		return nil
	}

	return setFieldValue(field, envValue, envVarName)
}

// setFieldValue sets the field value based on its type
func setFieldValue(field reflect.Value, value string, envVarName string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Type() == reflect.TypeOf(time.Duration(0)) {
			duration, err := time.ParseDuration(value)
			if err != nil {
				return fmt.Errorf("invalid duration value for %s: %s", envVarName, value)
			}
			field.SetInt(int64(duration))
		} else {
			intVal, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid integer value for %s: %s", envVarName, value)
			}
			field.SetInt(intVal)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid unsigned integer value for %s: %s", envVarName, value)
		}
		field.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid float value for %s: %s", envVarName, value)
		}
		field.SetFloat(floatVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid boolean value for %s: %s", envVarName, value)
		}
		field.SetBool(boolVal)
	default:
		return fmt.Errorf("unsupported field type %s for %s", field.Kind(), envVarName)
	}

	return nil
}

// buildEnvVarName builds the environment variable name from prefix and field path
func buildEnvVarName(prefix, parentPath, fieldName string) string {
	parts := []string{CamelToScreamingSnake(prefix)}

	if parentPath != "" {
		parts = append(parts, CamelToScreamingSnake(parentPath))
	}

	parts = append(parts, CamelToScreamingSnake(fieldName))

	return strings.Join(parts, "_")
}

// buildPath builds the path for nested structs
func buildPath(parentPath, fieldName string) string {
	if parentPath == "" {
		return fieldName
	}

	return parentPath + "_" + fieldName
}

// CamelToScreamingSnake converts camelCase to SCREAMING_SNAKE_CASE
func CamelToScreamingSnake(s string) string {
	var result strings.Builder

	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			// Check if previous character was lowercase or if next character is lowercase
			prevChar := rune(s[i-1])
			nextIsLower := i+1 < len(s) && unicode.IsLower(rune(s[i+1]))

			if unicode.IsLower(prevChar) || nextIsLower {
				result.WriteRune('_')
			}
		}
		result.WriteRune(unicode.ToLower(r))
	}

	return strings.ToUpper(result.String())
}
