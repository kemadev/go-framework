/*
Copyright 2025 kemadev
SPDX-License-Identifier: MPL-2.0
*/

package main

import (
	"log/slog"
	"os"

	"github.com/kemadev/go-framework/pkg/config"
	"github.com/kemadev/go-framework/pkg/log"
	"github.com/kemadev/go-framework/pkg/route"
	"github.com/kemadev/go-framework/pkg/serve"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

func main() {
	// Get app config
	conf, err := config.Load()
	if err != nil {
		log.CreateFallbackLogger(config.Runtime{}).Error(
			"run",
			slog.String("Body", "config failure"),
			// TODO use semconv value once released, see https://opentelemetry.io/docs/specs/semconv/attributes-registry/error/#error-message
			slog.String(string(semconv.ErrorMessageKey), err.Error()),
		)
		os.Exit(1)
	}

	// `http.Run()` only returns on init / shutdown failures, where otel logger isn't available
	fallbackLogger := log.CreateFallbackLogger(conf.Runtime)

	// Create regular routes
	regularRoutes := route.RoutesToRegister{
		route.Route{
			Pattern:     "GET /rolldice/",
			HandlerFunc: rolldice,
		},
		route.Route{
			Pattern:     "GET /rolldice/{player}",
			HandlerFunc: rolldice,
		},
	}

	// Create routes with dependency injection
	dependencyRoutes := route.RoutesWithDependencies{}

	// Run HTTP server
	err = serve.Run(regularRoutes, dependencyRoutes, conf)
	if err != nil {
		fallbackLogger.Error(
			"run",
			slog.String("Body", "http failure"),
			// TODO use semconv value once released, see https://opentelemetry.io/docs/specs/semconv/attributes-registry/error/#error-message
			slog.String("error.message", err.Error()),
		)
		os.Exit(1)
	}
}
