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

// Global is the server configuration struct.
// Values are populated from environment variables nammed after
// their relative position in the struct with "kema" as prefix, using SCREAMING_SNAKE_CASE.
// e.g. [Global.Observability.EndpointURL] is pupulated from environment variable `KEMA_OBSERVABILITY_ENDPOINT_URL`
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
	ListenAddr string `default:"[::]"`
	// Server bind port
	ListenPort int `default:"8080"`
	// HTTP read timeout
	ReadTimeout time.Duration `default:"15s"`
	// HTTP write timeout
	WriteTimeout time.Duration `default:"15s"`
	// HTTP idle timeout
	IdleTimeout time.Duration `default:"60s"`
	// Proxy header for forwarded entity
	ProxyHeader string `default:"X-Forwarded-For"`
}

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

type Observability struct {
	// Address of OpenTelemetry endpoint where to send telemetry
	EndpointURL string `required:"true"`
	// Compression to use when sending telemetry
	ExporterCompression string `default:"gzip"`
	// Percentage of request to sample for tracing
	TracingSamplePercent int `default:"100"`
	// Interval between metrics exports, in seconds
	MetricsExportIntervalSeconds int `default:"15"`
}

// Load loads configuration from environment variables
func Load() (*Global, error) {
	var cfg Global
	err := load("kema", &cfg)
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

		// Build the environment variable name
		fieldName := fieldType.Name
		envVarName := buildEnvVarName(prefix, parentPath, fieldName)

		switch field.Kind() {
		case reflect.Struct:
			// Recursively process nested structs
			err := processStruct(prefix, field, buildPath(parentPath, fieldName))
			if err != nil {
				return fmt.Errorf("error processing struct field %s: %w", fieldName, err)
			}
		default:
			// Process primitive fields
			err := processField(field, fieldType, envVarName)
			if err != nil {
				return fmt.Errorf("error processing field %s: %w", fieldType.Name, err)
			}
		}
	}

	return nil
}

// processField processes a single struct field
func processField(field reflect.Value, fieldType reflect.StructField, envVarName string) error {
	// Get tags
	defaultValue := fieldType.Tag.Get("default")
	required := fieldType.Tag.Get("required") == "true"

	// Get environment variable value
	envValue := os.Getenv(envVarName)

	// Use default if no env value and default is provided
	if envValue == "" && defaultValue != "" {
		envValue = defaultValue
	}

	// Check required fields
	if required && envValue == "" {
		return fmt.Errorf("required environment variable %s is not set", envVarName)
	}

	// Skip if no value and not required
	if envValue == "" {
		return nil
	}

	// Set the field value based on its type
	return setFieldValue(field, envValue, envVarName)
}

// setFieldValue sets the field value based on its type
func setFieldValue(field reflect.Value, value string, envVarName string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Type() == reflect.TypeOf(time.Duration(0)) {
			// Handle time.Duration
			duration, err := time.ParseDuration(value)
			if err != nil {
				return fmt.Errorf("invalid duration value for %s: %s", envVarName, value)
			}
			field.SetInt(int64(duration))
		} else {
			// Handle regular integers
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
	parts := []string{prefix}

	if parentPath != "" {
		parts = append(parts, parentPath)
	}

	// Convert camelCase to snake_case
	snakeCase := CamelToScreamingSnake(fieldName)
	parts = append(parts, snakeCase)

	// Output as SCREAMING_SNAKE
	return strings.ToUpper(strings.Join(parts, "_"))
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
