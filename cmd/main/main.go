/*
Copyright 2025 kemadev
SPDX-License-Identifier: MPL-2.0
*/

package main

import (
	"context"
	"errors"
	"fmt"
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
	ctx, stopSig := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stopSig()

	// Get app config
	conf, err := config.Load()

	// Set up OpenTelemetry.
	otelShutdown, err := otel.SetupOTelSDK(ctx, *conf)
	if err != nil {
		log.FallbackError(fmt.Errorf("failure setting up OpenTelemetry SDK: %w", err))
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
		ReadTimeout:  conf.Server.ReadTimeout * time.Second,
		WriteTimeout: conf.Server.WriteTimeout * time.Second,
		IdleTimeout:  conf.Server.IdleTimeout * time.Second,
		ErrorLog: slog.NewLogLogger(
			otelslog.NewLogger("net/http").Handler(),
			func() slog.Level {
				if conf.Runtime.IsLocalEnvironment() {
					return slog.LevelDebug
				}
				return slog.LevelError
			}(),
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
		log.FallbackError(fmt.Errorf("failure running HTTP server: %w", err))
		os.Exit(1)
	case <-ctx.Done():
		// Stop receiving signal notifications as soon as possible.
		stopSig()
	}

	err = srv.Shutdown(context.Background())
	if err != nil {
		log.FallbackError(fmt.Errorf("failure shutting down HTTP server: %w", err))
		os.Exit(1)
	}
}
