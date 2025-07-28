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
	sigCtx, stopSig := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stopSig()

	// Get app config
	conf, err := config.Load()
	if err != nil {
		log.FallbackError(fmt.Errorf("failure loading config: %w", err))
		os.Exit(1)
	}

	// Set up OpenTelemetry.
	otelShutdown, err := otel.SetupOTelSDK(sigCtx, *conf)
	if err != nil {
		log.FallbackError(fmt.Errorf("failure setting up OpenTelemetry SDK: %w", err))
		os.Exit(1)
	}

	// Global program return code
	var exitCode int

	// Cleanup function
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			conf.Observability.ShutdownGracePeriod*time.Second,
		)
		defer cancel()

		shutdownErr := otelShutdown(shutdownCtx)
		if shutdownErr != nil {
			log.FallbackError(fmt.Errorf("failure shutting down OpenTelemetry: %w", shutdownErr))
			if exitCode == 0 {
				exitCode = 1
			}
		}

		if exitCode != 0 {
			os.Exit(exitCode)
		}
	}()

	// Start HTTP server.
	srv := &http.Server{
		Addr:         conf.Server.BindAddr + ":" + strconv.Itoa(conf.Server.BindPort),
		BaseContext:  func(_ net.Listener) context.Context { return sigCtx },
		ReadTimeout:  conf.Server.ReadTimeout,
		WriteTimeout: conf.Server.WriteTimeout,
		IdleTimeout:  conf.Server.IdleTimeout,
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
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.FallbackError(fmt.Errorf("failure running HTTP server: %w", err))
			exitCode = 1
			return
		}
	case <-sigCtx.Done():
		// Stop receiving signal notifications as soon as possible.
		stopSig()
	}

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		// Let connections close, plus a grace period
		func() time.Duration {
			return max(
				conf.Server.ReadTimeout,
				conf.Server.WriteTimeout,
			) + conf.Server.ShutdownGracePeriod*time.Second
		}(),
	)
	defer cancel()

	err = srv.Shutdown(shutdownCtx)
	if err != nil {
		log.FallbackError(fmt.Errorf("failure shutting down HTTP server: %w", err))
		exitCode = 1
	}
}
