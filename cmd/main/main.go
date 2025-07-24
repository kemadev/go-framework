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
)

func main() {
	// `http.Run()` only returns on init / shutdown failures, where otel logger isn't available
	fallbackLogger := log.CreateFallbackLogger()

	// Get app config
	conf, err := config.NewConfig()
	if err != nil {
		fallbackLogger.Error(
			"run",
			slog.String("Body", "config failure"),
			// TODO use semconv value once released, see https://opentelemetry.io/docs/specs/semconv/attributes-registry/error/#error-message
			slog.String("error.message", err.Error()),
		)
		os.Exit(1)
	}

	// Define routes to handle
	routes := serve.RoutesToRegister{
		route.Route{
			Pattern:     "/rolldice/",
			HandlerFunc: rolldice,
		},
		route.Route{
			Pattern:     "/rolldice/{player}",
			HandlerFunc: rolldice,
		},
	}

	// Run HTTP server
	err = serve.Run(routes, conf)
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
