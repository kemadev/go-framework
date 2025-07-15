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
		otelEnpointURL          string
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
			runtimeEnvVarValue:      config.EnvMain,
			appName:                 "app",
			appNamespace:            "namespace",
			otelEnpointURL:          "grpc://localhost:4317",
			otelExporterCompression: "gzip",
			servePort:               "8080",
			writeTimeout:            "10",
			readTimeout:             "1",
			appVersion:              "v0.0.1",
			want: &config.Config{
				RuntimeEnv:              config.EnvMain,
				OtelEndpointURL:         url.URL{Scheme: "grpc", Host: "localhost:4317"},
				AppVersion:              "v0.0.1",
				AppName:                 "app",
				AppNamespace:            "namespace",
				OtelExporterCompression: "gzip",
				HTTPServePort:           8080,
				HTTPWriteTimeout:        10,
				HTTPReadTimeout:         1,
			},
			wantErr: false,
		},
		{
			name:                    "invalid runtime env",
			runtimeEnvVarValue:      "invalid",
			appName:                 "app",
			appNamespace:            "namespace",
			otelEnpointURL:          "grpc://localhost:4317",
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
			otelEnpointURL:          "grpc://localhost:4317",
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
			runtimeEnvVarValue:      config.EnvMain,
			appName:                 "app",
			appNamespace:            "namespace",
			otelEnpointURL:          "://bad-url",
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
			runtimeEnvVarValue:      config.EnvMain,
			appName:                 "app",
			appNamespace:            "namespace",
			otelEnpointURL:          "grpc://localhost:4317",
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
			runtimeEnvVarValue:      config.EnvMain,
			appName:                 "app",
			appNamespace:            "namespace",
			otelEnpointURL:          "grpc://localhost:4317",
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
			runtimeEnvVarValue:      config.EnvMain,
			appName:                 "",
			appNamespace:            "namespace",
			otelEnpointURL:          "grpc://localhost:4317",
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
			runtimeEnvVarValue:      config.EnvMain,
			appName:                 "app",
			appNamespace:            "",
			otelEnpointURL:          "grpc://localhost:4317",
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
			runtimeEnvVarValue:      config.EnvMain,
			appName:                 "app",
			appNamespace:            "namespace",
			otelEnpointURL:          "grpc://localhost:4317",
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
			runtimeEnvVarValue:      config.EnvMain,
			appName:                 "app",
			appNamespace:            "namespace",
			otelEnpointURL:          "",
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
			runtimeEnvVarValue:      config.EnvMain,
			appName:                 "app",
			appNamespace:            "namespace",
			otelEnpointURL:          "grpc://localhost:4317",
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
			runtimeEnvVarValue:      config.EnvMain,
			appName:                 "app",
			appNamespace:            "namespace",
			otelEnpointURL:          "grpc://localhost:4317",
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
			runtimeEnvVarValue:      config.EnvMain,
			appName:                 "app",
			appNamespace:            "namespace",
			otelEnpointURL:          "grpc://localhost:4317",
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
			runtimeEnvVarValue:      config.EnvMain,
			appName:                 "app",
			appNamespace:            "namespace",
			otelEnpointURL:          "grpc://localhost:4317",
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
			runtimeEnvVarValue:      config.EnvDev,
			appName:                 "devapp",
			appNamespace:            "devns",
			otelEnpointURL:          "grpc://devhost:4317",
			otelExporterCompression: "gzip",
			servePort:               "3000",
			writeTimeout:            "5",
			readTimeout:             "2",
			appVersion:              "1.2.3",
			want: &config.Config{
				RuntimeEnv:              config.EnvDev,
				OtelEndpointURL:         url.URL{Scheme: "grpc", Host: "devhost:4317"},
				AppVersion:              "1.2.3",
				AppName:                 "devapp",
				AppNamespace:            "devns",
				OtelExporterCompression: "gzip",
				HTTPServePort:           3000,
				HTTPWriteTimeout:        5,
				HTTPReadTimeout:         2,
			},
			wantErr: false,
		},
		{
			name:                    "valid config - next environment",
			runtimeEnvVarValue:      config.EnvNext,
			appName:                 "nextapp",
			appNamespace:            "nextns",
			otelEnpointURL:          "grpc://nexthost:4317",
			otelExporterCompression: "gzip",
			servePort:               "4000",
			writeTimeout:            "15",
			readTimeout:             "3",
			appVersion:              "2.0.0",
			want: &config.Config{
				RuntimeEnv:              config.EnvNext,
				OtelEndpointURL:         url.URL{Scheme: "grpc", Host: "nexthost:4317"},
				AppVersion:              "2.0.0",
				AppName:                 "nextapp",
				AppNamespace:            "nextns",
				OtelExporterCompression: "gzip",
				HTTPServePort:           4000,
				HTTPWriteTimeout:        15,
				HTTPReadTimeout:         3,
			},
			wantErr: false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			envVarsMappings := map[string]string{
				testCase.runtimeEnvVarValue:      config.EnvVarKeyRuntimeEnv,
				testCase.otelEnpointURL:          config.EnvVarKeyOtelEndpointURL,
				testCase.otelExporterCompression: config.EnvVarKeyOtelExporterCompression,
				testCase.appVersion:              config.EnvVarKeyAppVersion,
				testCase.appName:                 config.EnvVarKeyAppName,
				testCase.appNamespace:            config.EnvVarKeyAppNamespace,
				testCase.servePort:               config.EnvVarKeyHTTPServePort,
				testCase.readTimeout:             config.EnvVarKeyHTTPReadTimeout,
				testCase.writeTimeout:            config.EnvVarKeyHTTPWriteTimeout,
			}
			for k, v := range envVarsMappings {
				if k != "" {
					t.Setenv(v, k)
				}
			}

			got, err := config.NewConfig()
			if (err != nil) != testCase.wantErr {
				t.Errorf("error = %v, wantErr %v", err, testCase.wantErr)

				return
			}

			if !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("NewConfig() = %v, want %v", got, testCase.want)
			}
		})
	}
}
