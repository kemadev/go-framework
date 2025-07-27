package server

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/contrib/otelfiber/v3"
	"github.com/gofiber/fiber/v3"
	"github.com/kemadev/go-framework/pkg/config"
	"github.com/kemadev/go-framework/pkg/log"
	"github.com/kemadev/go-framework/pkg/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

func New() {
	// Configure server
	conf, err := config.Load()
	if err != nil {
		slog.New(slog.NewJSONHandler(os.Stdout, nil)).
			Error("can't load config", slog.String("error.message", err.Error()))
		os.Exit(1)
	}

	app := fiber.New(fiber.Config{
		AppName:      conf.Runtime.AppName,
		ReadTimeout:  conf.Server.ReadTimeout,
		WriteTimeout: conf.Server.WriteTimeout,
		IdleTimeout:  conf.Server.IdleTimeout,
		ProxyHeader:  conf.Server.ProxyHeader,
	})

	// Set up OpenTelemetry.
	otelShutdown, err := otel.SetupOTelSDK(context.TODO(), *conf)
	if err != nil {
		log.CreateFallbackLogger(conf.Runtime).
			Error("error setting up OpenTelemetry", slog.String(string(semconv.ErrorMessageKey), err.Error()))
		return
	}

	app.Use(otelfiber.Middleware(
		otelfiber.WithCollectClientIP(true),
		otelfiber.WithPort(conf.Server.ListenPort),
	))

	app.Get("/foo/:user", func(c fiber.Ctx) error {
		// Send a string response to the client
		return c.SendString("Hello, World 👋!")
	})

	srvErr := make(chan error, 1)

	go func() {
		err := app.Listen(
			conf.Server.ListenAddr+":"+strconv.Itoa(conf.Server.ListenPort),
			fiber.ListenConfig{
				EnablePrefork:         !conf.Runtime.IsLocalEnvironment(),
				DisableStartupMessage: !conf.Runtime.IsLocalEnvironment(),
				EnablePrintRoutes:     conf.Runtime.IsLocalEnvironment(),
				ShutdownTimeout: func() time.Duration {
					return max(
						conf.Server.ReadTimeout,
						conf.Server.WriteTimeout,
						conf.Server.IdleTimeout,
					)
				}(),
				ListenerNetwork: func() string {
					if strings.HasPrefix(conf.Server.ListenAddr, "[") {
						return "tcp6"
					}
					return "tcp4"
				}(),
			},
		)
		if err != nil {
			srvErr <- err
		}
	}()

	// Set up graceful shutdown
	app.Hooks().OnPostShutdown(func(err error) error {
		logger := log.CreateFallbackLogger(conf.Runtime)
		if err != nil {
			logger.Error(
				"error shutting down the server",
				slog.String(string(semconv.ErrorMessageKey), err.Error()),
			)
		}
		otelErr := errors.Join(err, otelShutdown(context.Background()))
		if otelErr != nil {
			logger.Error(
				"error shutting down open telemetry",
				slog.String(string(semconv.ErrorMessageKey), otelErr.Error()),
			)
		}
		return nil
	})

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case err := <-srvErr:
		log.CreateFallbackLogger(conf.Runtime).
			Error("error starting server", slog.String(string(semconv.ErrorMessageKey), err.Error()))
		os.Exit(1)
	case <-ctx.Done():
		stop()
		log.CreateFallbackLogger(conf.Runtime).
			Info("received shutdown signal, shutting down server")
		err := app.Shutdown()
		if err != nil {
			log.CreateFallbackLogger(conf.Runtime).
				Error("error shutting down server", slog.String(string(semconv.ErrorMessageKey), err.Error()))
			os.Exit(1)
		}
	}
}
