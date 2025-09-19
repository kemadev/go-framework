/*
Copyright 2025 kemadev
SPDX-License-Identifier: MPL-2.0
*/

package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/kemadev/go-framework/pkg/client/cache"
	"github.com/kemadev/go-framework/pkg/config"
	"github.com/kemadev/go-framework/pkg/convenience/headval"
	"github.com/kemadev/go-framework/pkg/convenience/log"
	"github.com/kemadev/go-framework/pkg/convenience/otel"
	"github.com/kemadev/go-framework/pkg/convenience/render"
	"github.com/kemadev/go-framework/pkg/convenience/resp"
	"github.com/kemadev/go-framework/pkg/convenience/sechead"
	"github.com/kemadev/go-framework/pkg/convenience/trace"
	"github.com/kemadev/go-framework/pkg/encoding"
	flog "github.com/kemadev/go-framework/pkg/log"
	"github.com/kemadev/go-framework/pkg/maxbytes"
	"github.com/kemadev/go-framework/pkg/monitoring"
	"github.com/kemadev/go-framework/pkg/router"
	"github.com/kemadev/go-framework/pkg/server"
	"github.com/kemadev/go-framework/pkg/timeout"
	"github.com/kemadev/go-framework/web"
	"github.com/valkey-io/valkey-go"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

const packageName = "github/go-framework/go-framework/cmd/main"

func main() {
	// Get app config
	configManager := config.NewManager()

	conf, err := configManager.Get()
	if err != nil {
		flog.FallbackError(fmt.Errorf("error getting config: %w", err))
		os.Exit(1)
	}

	// Create clients, for use in handlers
	client, err := cache.NewClient(conf.Client.Cache)
	if err != nil {
		flog.FallbackError(err)
		os.Exit(1)
	}
	defer client.Close()

	r := router.New()

	// Always protect your routes (you can further customize at handler / group level)
	r.Use(timeout.NewMiddleware(5 * time.Second))
	r.Use(maxbytes.NewMiddleware(100000))

	// Add other middlewares
	r.Use(encoding.DecompressMiddleware)
	r.Use(encoding.CompressMiddleware)

	// Add monitoring endpoints
	r.Handle(
		monitoring.LivenessHandler(
			// Add your check function
			func() monitoring.CheckResults { return monitoring.CheckResults{} },
		),
	)
	r.Handle(
		monitoring.ReadinessHandler(
			// Add your check function
			func() monitoring.CheckResults { return monitoring.CheckResults{} },
		),
	)

	// Add handlers
	r.Handle(
		otel.WrapHandler("GET /foo/{bar}", http.HandlerFunc(ExampleHandler)),
	)

	r.Handle(
		otel.WrapHandler(
			"GET /db",
			ExampleDBHandler(client),
		),
	)

	// Create groups (sub-groups are also possible)
	r.Group(func(r *router.Router) {
		// Secure frontend with security headers
		r.Use(sechead.NewMiddleware(sechead.SecHeadersDefaultStrict))

		// Handle template assets
		tmplFS := web.GetTmplFS()
		renderer, _ := render.New(tmplFS)
		r.Handle(
			otel.WrapHandler(
				"GET /",
				ExampleTemplateRender(renderer),
			),
		)
	})

	// Handle static (public) assets
	r.Handle(
		otel.WrapHandler(
			"GET /static/",
			http.FileServerFS(web.GetStaticFS()).ServeHTTP,
		),
	)

	server.Run(otel.WrapMux(r, packageName), *conf)
}

func ExampleHandler(w http.ResponseWriter, r *http.Request) {
	span := trace.Span(r.Context())
	span.SetAttributes(attribute.String("bar", r.PathValue("bar")))

	// Use otelhttp to call external services so it is automatically instrumented
	eresp, err := otelhttp.Get(r.Context(), "https://example.com")
	if err != nil {
		log.ErrLog(packageName, "error calling external http endpoint", err)
		// Prefer graceful degradation instead of throwing a 5XX error
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)

		return
	}

	type exampleResp struct {
		Name  string
		Attrs []string
	}

	var name []byte
	_, err = eresp.Body.Read(name)
	if err != nil {
		log.ErrLog(packageName, "error calling external http endpoint", err)
		// Prefer graceful degradation instead of throwing a 5XX error
		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)

		return
	}

	resp.JSON(w, exampleResp{
		Name:  string(name),
		Attrs: []string{"whatever"},
	})
}

func ExampleTemplateRender(
	tr *render.TemplateRenderer,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := tr.Execute(
			w,
			// Mind directory name in tmplate FS
			"tmpl/hello.html",
			map[string]any{
				"WorldName": "WoRlD",
			},
			headval.MIMETextHTMLCharsetUTF8,
		)
		if err != nil {
			if errors.Is(err, render.ErrTemplateNotFound) {
				http.NotFound(w, r)
				return
			}

			log.Logger(packageName).
				Error("error rendering template",
					slog.String(
						string(semconv.ErrorMessageKey),
						err.Error(),
					),
				)
			// Prefer graceful degradation instead of throwing a 5XX error
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
		}
	}
}

func ExampleDBHandler(client valkey.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := client.Do(r.Context(), client.B().Set().Key("key").Value(time.Now().String()).Build()).
			Error()
		if err != nil {
			log.ErrLog(packageName, "error db set", err)
			// Prefer graceful degradation instead of throwing a 5XX error
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)

			return
		}

		resp.JSON(w, "ok")
	}
}
