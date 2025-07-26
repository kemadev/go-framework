package server

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/gofiber/contrib/otelfiber/v3"
	"github.com/gofiber/fiber/v3"
	"github.com/kemadev/go-framework/pkg/config"
	"github.com/kemadev/go-framework/pkg/otel"
)

func New() {
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

	otel.SetupOTelSDK(context.TODO(), *conf)

	app.Use(otelfiber.Middleware())

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	err = app.Listen(conf.Server.ListenAddr, fiber.ListenConfig{
		EnablePrefork:         !conf.IsLocalEnvironment(),
		DisableStartupMessage: !conf.IsLocalEnvironment(),
		EnablePrintRoutes:     conf.IsLocalEnvironment(),
		ShutdownTimeout: func() time.Duration {
			return max(conf.Server.ReadTimeout, conf.Server.WriteTimeout, conf.Server.IdleTimeout)
		}(),
		ListenerNetwork: func() string {
			if strings.HasPrefix(conf.Server.ListenAddr, "[") {
				return "tcp6"
			}
			return "tcp4"
		}(),
	})
	if err != nil {
		slog.Error("server listen failed: %w", slog.String("error.message", err.Error()))
		os.Exit(1)
	}
}
