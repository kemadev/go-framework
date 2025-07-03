package otel

import (
	"context"
	"net/url"
	"testing"

	"github.com/kemadev/go-framework/pkg/config"
)

// TestSetupOTelSDK tests the SetupOTelSDK function.
func TestSetupOTelSDK(t *testing.T) {
	type args struct {
		ctx  context.Context
		conf config.Config
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "runtime-env-dev",
			args: args{
				ctx: t.Context(),
				conf: config.Config{
					RuntimeEnv: config.Env_dev,
					OtelEndpointUrl: url.URL{
						Scheme: "grpc",
						Host:   "localhost:4317",
					},
					OtelExporterCompression: "gzip",
				},
			},
			wantErr: false,
		},
		{
			name: "runtime-env-next",
			args: args{
				ctx: t.Context(),
				conf: config.Config{
					RuntimeEnv: config.Env_next,
					OtelEndpointUrl: url.URL{
						Scheme: "grpc",
						Host:   "localhost:4317",
					},
					OtelExporterCompression: "gzip",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SetupOTelSDK(tt.args.ctx, tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetupOTelSDK() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
		})
	}
}
