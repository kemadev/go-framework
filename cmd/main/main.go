/*
Copyright 2025 kemadev
SPDX-License-Identifier: MPL-2.0
*/

package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/kemadev/go-framework/pkg/config"
	"github.com/kemadev/go-framework/pkg/log"
	"github.com/kemadev/go-framework/pkg/otel"
	"go.opentelemetry.io/contrib/bridges/otelslog"
)

func main() {
	// Intercept signals
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	// Get app config
	conf, err := config.Load()
	if err != nil {
		log.FallbackError("error loading config", err)
		os.Exit(1)
	}

	// Set up OpenTelemetry.
	otelShutdown, err := otel.SetupOTelSDK(ctx, *conf)
	if err != nil {
		log.FallbackError("error setting up OpenTelemetry", err)
		os.Exit(1)
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// Start HTTP server.
	srv := &http.Server{
		// Use any host, let Kubernetes handle the routing.
		Addr:         ":" + strconv.Itoa(conf.Server.BindPort),
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Duration(conf.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(conf.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(conf.Server.IdleTimeout) * time.Second,
		ErrorLog: slog.NewLogLogger(
			otelslog.NewLogger("httpserver").Handler(),
			slog.LevelError,
		),
		Handler: nil,
	}

	srvErr := make(chan error, 1)

	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	// Wait for interruption
	select {
	case err = <-srvErr:
		log.FallbackError("HTTP server error", err)
		os.Exit(1)
	case <-ctx.Done():
		// Stop receiving signal notifications as soon as possible.
		stop()
	}

	err = srv.Shutdown(context.Background())
	if err != nil {
		log.FallbackError("HTTP server shutdown error", err)
		os.Exit(1)
	}
}
