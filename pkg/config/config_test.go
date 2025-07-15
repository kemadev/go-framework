// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package config_test

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/kemadev/go-framework/pkg/config"
)

// TestNewConfig tests the NewConfig function from the config package.
func TestNewConfig(t *testing.T) {
	tests := []struct {
		name                    string
		runtimeEnvVarValue      string
		appName                 string
		appNamespace            string
		otelEnpointUrl          string
		otelExporterCompression string
		servePort               string
		writeTimeout            string
		readTimeout             string
		appVersion              string
		want                    *config.Config
		wantErr                 bool
	}{
		{
			name:                    "valid config - all fields set",
			runtimeEnvVarValue:      config.Env_main,
			appName:                 "app",
			appNamespace:            "namespace",
			otelEnpointUrl:          "grpc://localhost:4317",
			otelExporterCompression: "gzip",
			servePort:               "8080",
			writeTimeout:            "10",
			readTimeout:             "1",
			appVersion:              "v0.0.1",
			want: &config.Config{
				RuntimeEnv:              config.Env_main,
				OtelEndpointUrl:         url.URL{Scheme: "grpc", Host: "localhost:4317"},
				AppVersion:              "v0.0.1",
				AppName:                 "app",
				AppNamespace:            "namespace",
				OtelExporterCompression: "gzip",
				HttpServePort:           8080,
				HttpWriteTimeout:        10,
				HttpReadTimeout:         1,
			},
			wantErr: false,
		},
		{
			name:                    "invalid runtime env",
			runtimeEnvVarValue:      "invalid",
			appName:                 "app",
			appNamespace:            "namespace",
			otelEnpointUrl:          "grpc://localhost:4317",
			otelExporterCompression: "gzip",
			servePort:               "8080",
			writeTimeout:            "10",
			readTimeout:             "1",
			appVersion:              "v0.0.1",
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "missing required env var",
			runtimeEnvVarValue:      "",
			appName:                 "app",
			appNamespace:            "namespace",
			otelEnpointUrl:          "grpc://localhost:4317",
			otelExporterCompression: "gzip",
			servePort:               "8080",
			writeTimeout:            "10",
			readTimeout:             "1",
			appVersion:              "v0.0.1",
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "invalid otel endpoint url",
			runtimeEnvVarValue:      config.Env_main,
			appName:                 "app",
			appNamespace:            "namespace",
			otelEnpointUrl:          "://bad-url",
			otelExporterCompression: "gzip",
			servePort:               "8080",
			writeTimeout:            "10",
			readTimeout:             "1",
			appVersion:              "v0.0.1",
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "invalid serve port",
			runtimeEnvVarValue:      config.Env_main,
			appName:                 "app",
			appNamespace:            "namespace",
			otelEnpointUrl:          "grpc://localhost:4317",
			otelExporterCompression: "gzip",
			servePort:               "notaport",
			writeTimeout:            "10",
			readTimeout:             "1",
			appVersion:              "v0.0.1",
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "invalid otel compression",
			runtimeEnvVarValue:      config.Env_main,
			appName:                 "app",
			appNamespace:            "namespace",
			otelEnpointUrl:          "grpc://localhost:4317",
			otelExporterCompression: "snappy",
			servePort:               "8080",
			writeTimeout:            "10",
			readTimeout:             "1",
			appVersion:              "v0.0.1",
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "empty app name",
			runtimeEnvVarValue:      config.Env_main,
			appName:                 "",
			appNamespace:            "namespace",
			otelEnpointUrl:          "grpc://localhost:4317",
			otelExporterCompression: "gzip",
			servePort:               "8080",
			writeTimeout:            "10",
			readTimeout:             "1",
			appVersion:              "v0.0.1",
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "empty app namespace",
			runtimeEnvVarValue:      config.Env_main,
			appName:                 "app",
			appNamespace:            "",
			otelEnpointUrl:          "grpc://localhost:4317",
			otelExporterCompression: "gzip",
			servePort:               "8080",
			writeTimeout:            "10",
			readTimeout:             "1",
			appVersion:              "v0.0.1",
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "empty app version",
			runtimeEnvVarValue:      config.Env_main,
			appName:                 "app",
			appNamespace:            "namespace",
			otelEnpointUrl:          "grpc://localhost:4317",
			otelExporterCompression: "gzip",
			servePort:               "8080",
			writeTimeout:            "10",
			readTimeout:             "1",
			appVersion:              "",
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "empty otel endpoint url",
			runtimeEnvVarValue:      config.Env_main,
			appName:                 "app",
			appNamespace:            "namespace",
			otelEnpointUrl:          "",
			otelExporterCompression: "gzip",
			servePort:               "8080",
			writeTimeout:            "10",
			readTimeout:             "1",
			appVersion:              "v0.0.1",
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "empty otel exporter compression",
			runtimeEnvVarValue:      config.Env_main,
			appName:                 "app",
			appNamespace:            "namespace",
			otelEnpointUrl:          "grpc://localhost:4317",
			otelExporterCompression: "",
			servePort:               "8080",
			writeTimeout:            "10",
			readTimeout:             "1",
			appVersion:              "v0.0.1",
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "empty serve port",
			runtimeEnvVarValue:      config.Env_main,
			appName:                 "app",
			appNamespace:            "namespace",
			otelEnpointUrl:          "grpc://localhost:4317",
			otelExporterCompression: "gzip",
			servePort:               "",
			writeTimeout:            "10",
			readTimeout:             "1",
			appVersion:              "v0.0.1",
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "empty write timeout",
			runtimeEnvVarValue:      config.Env_main,
			appName:                 "app",
			appNamespace:            "namespace",
			otelEnpointUrl:          "grpc://localhost:4317",
			otelExporterCompression: "gzip",
			servePort:               "8080",
			writeTimeout:            "",
			readTimeout:             "1",
			appVersion:              "v0.0.1",
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "empty read timeout",
			runtimeEnvVarValue:      config.Env_main,
			appName:                 "app",
			appNamespace:            "namespace",
			otelEnpointUrl:          "grpc://localhost:4317",
			otelExporterCompression: "gzip",
			servePort:               "8080",
			writeTimeout:            "10",
			readTimeout:             "",
			appVersion:              "v0.0.1",
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "valid config - dev environment",
			runtimeEnvVarValue:      config.Env_dev,
			appName:                 "devapp",
			appNamespace:            "devns",
			otelEnpointUrl:          "grpc://devhost:4317",
			otelExporterCompression: "gzip",
			servePort:               "3000",
			writeTimeout:            "5",
			readTimeout:             "2",
			appVersion:              "1.2.3",
			want: &config.Config{
				RuntimeEnv:              config.Env_dev,
				OtelEndpointUrl:         url.URL{Scheme: "grpc", Host: "devhost:4317"},
				AppVersion:              "1.2.3",
				AppName:                 "devapp",
				AppNamespace:            "devns",
				OtelExporterCompression: "gzip",
				HttpServePort:           3000,
				HttpWriteTimeout:        5,
				HttpReadTimeout:         2,
			},
			wantErr: false,
		},
		{
			name:                    "valid config - next environment",
			runtimeEnvVarValue:      config.Env_next,
			appName:                 "nextapp",
			appNamespace:            "nextns",
			otelEnpointUrl:          "grpc://nexthost:4317",
			otelExporterCompression: "gzip",
			servePort:               "4000",
			writeTimeout:            "15",
			readTimeout:             "3",
			appVersion:              "2.0.0",
			want: &config.Config{
				RuntimeEnv:              config.Env_next,
				OtelEndpointUrl:         url.URL{Scheme: "grpc", Host: "nexthost:4317"},
				AppVersion:              "2.0.0",
				AppName:                 "nextapp",
				AppNamespace:            "nextns",
				OtelExporterCompression: "gzip",
				HttpServePort:           4000,
				HttpWriteTimeout:        15,
				HttpReadTimeout:         3,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			envVarsMappings := map[string]string{
				tt.runtimeEnvVarValue:      config.EnvVarKeyRuntimeEnv,
				tt.otelEnpointUrl:          config.EnvVarKeyOtelEndpointURL,
				tt.otelExporterCompression: config.EnvVarKeyOtelExporterCompression,
				tt.appVersion:              config.EnvVarKeyAppVersion,
				tt.appName:                 config.EnvVarKeyAppName,
				tt.appNamespace:            config.EnvVarKeyAppNamespace,
				tt.servePort:               config.EnvVarKeyHTTPServePort,
				tt.readTimeout:             config.EnvVarKeyHTTPReadTimeout,
				tt.writeTimeout:            config.EnvVarKeyHTTPWriteTimeout,
			}
			for k, v := range envVarsMappings {
				if k != "" {
					t.Setenv(v, k)
				}
			}

			got, err := config.NewConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
