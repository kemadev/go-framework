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

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kemadev/go-framework/pkg/client/cache"
	"github.com/kemadev/go-framework/pkg/client/database"
	"github.com/kemadev/go-framework/pkg/client/search"
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
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/valkey-io/valkey-go"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

const packageName = "github.com/kemadev/go-framework/cmd/go-framework"

func main() {
	// Get app config
	configManager := config.NewManager()

	conf, err := configManager.Get()
	if err != nil {
		flog.FallbackError(fmt.Errorf("error getting config: %w", err))
		os.Exit(1)
	}

	// Create clients, for use in handlers
	cacheClient, err := cache.NewClient(conf.Client.Cache)
	if err != nil {
		flog.FallbackError(err)
		os.Exit(1)
	}
	defer cacheClient.Close()

	databaseClient, err := database.NewClient(conf.Client.Database)
	if err != nil {
		flog.FallbackError(err)
		os.Exit(1)
	}
	defer databaseClient.Close()

	searchClient, err := search.NewClient(conf.Client.Search, conf.Runtime)
	if err != nil {
		flog.FallbackError(err)
		os.Exit(1)
	}

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
			func() monitoring.CheckResults {
				// Add your check function logic
				return monitoring.CheckResults{}
			},
		),
	)
	r.Handle(
		monitoring.ReadinessHandler(
			func() monitoring.CheckResults {
				return monitoring.CheckResults{
					// Adjust status on ping fail
					"database": database.Check(databaseClient, monitoring.StatusDown),
					"cache":    cache.Check(cacheClient, monitoring.StatusDown),
					"search":   search.Check(searchClient, monitoring.StatusDown),
					// Add your check functions
				}
			},
		),
	)

	// Add handlers
	r.Handle(
		otel.WrapHandler("GET /foo/{bar}", http.HandlerFunc(ExampleHandler)),
	)

	r.Handle(
		otel.WrapHandler(
			"GET /cache",
			ExampleCacheHandler(cacheClient),
		),
	)

	r.Handle(
		otel.WrapHandler(
			"GET /database",
			ExampleDatabaseHandler(databaseClient),
		),
	)

	r.Handle(
		otel.WrapHandler(
			"GET /search",
			ExampleSearchHandler(searchClient),
		),
	)

	// Create groups (sub-groups are also possible)
	r.Group(func(r *router.Router) {
		// Secure frontend with security headers
		r.Use(sechead.NewMiddleware(sechead.SecHeadersDefaultStrict))

		// Handle template assets
		tmplFS := web.GetTmplFS()
		renderer, _ := render.New(tmplFS, web.TemplateBaseDirName)
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
			"GET /"+web.StaticBaseDirName+"/",
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
			// Mind about file extension
			r.URL.Path+".gotmpl.html",
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

func ExampleCacheHandler(client valkey.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := client.Do(r.Context(), client.B().Set().Key("key").Value(time.Now().String()).Build()).
			Error()
		if err != nil {
			log.ErrLog(packageName, "error cache set", err)
			// Prefer graceful degradation instead of throwing a 5XX error
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)

			return
		}

		type ExampleOutput struct {
			Success bool
		}

		resp.JSON(w, ExampleOutput{
			Success: true,
		})
	}
}

func ExampleDatabaseHandler(client *pgxpool.Pool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var id int

		err := client.QueryRow(
			r.Context(),
			`INSERT INTO tasks (created_at) VALUES ($1) RETURNING id`,
			time.Now(),
		).Scan(&id)
		if err != nil {
			log.ErrLog(packageName, "error database insert", err)
			// Prefer graceful degradation instead of throwing a 5XX error
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)

			return
		}

		type ExampleOutput struct {
			ID int
		}

		resp.JSON(w, ExampleOutput{ID: id})
	}
}

func ExampleSearchHandler(
	client *opensearchapi.Client,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := client.Info(r.Context(), nil)
		if err != nil {
			log.ErrLog(packageName, "error search info", err)
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
			return
		}

		type ExampleOutput struct {
			ClusterName string
		}

		resp.JSON(w, ExampleOutput{
			ClusterName: res.ClusterName,
		})
	}
}
