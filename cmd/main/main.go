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
			slog.String("Body", "error loading config"),
			slog.String(string(semconv.ErrorMessageKey), err.Error()),
		)
		os.Exit(1)
	}

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
		log.CreateFallbackLogger(conf.Runtime).Error(
			"run",
			slog.String("Body", "error running server"),
			slog.String(string(semconv.ErrorMessageKey), err.Error()),
		)
		os.Exit(1)
	}
}
