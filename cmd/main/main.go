/*
Copyright 2025 kemadev
SPDX-License-Identifier: MPL-2.0
*/

package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/kemadev/go-framework/pkg/client"
	"github.com/kemadev/go-framework/pkg/config"
	"github.com/kemadev/go-framework/pkg/convenience/headval"
	"github.com/kemadev/go-framework/pkg/convenience/local"
	"github.com/kemadev/go-framework/pkg/convenience/log"
	"github.com/kemadev/go-framework/pkg/convenience/otel"
	"github.com/kemadev/go-framework/pkg/convenience/render"
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
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const packageName = "github.com/kemadev/go-framework/cmd/main"

func main() {
	// Get app config
	configManager := config.NewManager()

	conf, err := configManager.Get()
	if err != nil {
		flog.FallbackError(fmt.Errorf("error getting config: %w", err))
		os.Exit(1)
	}

	// Create clients, for use in handlers
	client, err := client.NewValkeyDBClient(conf.Client.Database)
	if err != nil {
		flog.FallbackError(err)
		os.Exit(1)
	}
	defer client.Close()

	r := router.New()

	// Always protect your routes (you can further customize at handler / group level)
	r.Use(otel.WrapMiddleware("timeout", timeout.NewMiddleware(5*time.Second)))
	r.Use(otel.WrapMiddleware("maxbytes", maxbytes.NewMiddleware(100000)))

	// Add other middlewares
	r.Use(otel.WrapMiddleware("decompress", encoding.DecompressMiddleware))
	r.Use(otel.WrapMiddleware("compress", encoding.CompressMiddleware))
	r.Use(otel.WrapMiddleware("logging", LoggingMiddleware))

	// Add monitoring endpoints
	r.Handle(
		monitoring.LivenessHandler(
			func() monitoring.CheckResults { return monitoring.CheckResults{} },
		),
	)
	r.Handle(
		monitoring.ReadinessHandler(
			func() monitoring.CheckResults { return monitoring.CheckResults{} },
		),
	)

	// Add handlers
	r.Handle(
		otel.WrapHandler("GET /foo/{bar}", http.HandlerFunc(FooBar)),
	)
	r.Handle(
		otel.WrapHandler(
			"POST /tester",
			http.HandlerFunc(Tester),
		),
	)

	// Create groups
	r.Group(func(r *router.Router) {
		r.Use(otel.WrapMiddleware("auth", AuthMiddleware))

		r.Group(func(r *router.Router) {
			r.Use(otel.WrapMiddleware("timing", TimingMiddleware))

			r.Handle(
				otel.WrapHandler(
					"GET /auth/{bar}",
					http.HandlerFunc(FooBar),
				),
			)
		})
	})

	// Handle static (public) assets
	r.Handle(
		otel.WrapHandler(
			"GET /static/",
			http.FileServerFS(web.GetStaticFS()).ServeHTTP,
		),
	)

	// Handle template assets
	tmplFS := web.GetTmplFS()
	renderer, _ := render.New(tmplFS)
	r.Handle(
		otel.WrapHandler(
			"GET /",
			ExampleTemplateRender(renderer),
		),
	)

	r.Handle(
		otel.WrapHandler(
			"GET /sleep",
			timeout.WrapHandler(http.HandlerFunc(Wait), 1*time.Second).ServeHTTP,
		),
	)

	r.Handle(
		otel.WrapHandler(
			"GET /db",
			DBConn(client),
		),
	)

	server.Run(otel.WrapMux(r, packageName), *conf)
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
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func FooBar(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by AuthMiddleware)
	user := local.Get(r.Context(), "user")

	// Get span context for logging
	span := trace.Span(r.Context())
	spanCtx := span.SpanContext()
	fmt.Printf("[HANDLER] TraceID: %s, SpanID: %s, User: %v\n",
		spanCtx.TraceID().String(),
		spanCtx.SpanID().String(),
		user,
	)

	bag := trace.Baggage(r.Context())
	span.AddEvent(
		"handling this...",
		oteltrace.WithAttributes(semconv.UserID(bag.Member(string(semconv.UserIDKey)).Value())),
	)

	span.SetAttributes(attribute.String("bar", r.PathValue("bar")))

	fmt.Fprintf(w, "Hello, %v! TraceID: %s", user, spanCtx.TraceID().String())
}

type TesterPayload2 struct {
	Foo    string
	hidden string
}

type TesterPayload struct {
	hidden string
	Foo    string `json:"hello"`
	Quux   TesterPayload2
}

func Tester(w http.ResponseWriter, r *http.Request) {
	bod, err := io.ReadAll(r.Body)
	if err != nil {
		log.ErrLog("Tester", "read err", err)
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			http.Error(
				w,
				http.StatusText(http.StatusRequestEntityTooLarge),
				http.StatusRequestEntityTooLarge,
			)
			return
		}
	}
	slog.Debug(fmt.Sprintf("%s", bod))
	w.WriteHeader(200)
	w.Write([]byte("this is the response"))
}

func DBConn(client valkey.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := client.Do(r.Context(), client.B().Set().Key("key").Value(time.Now().String()).Build()).
			Error()
		if err != nil {
			log.ErrLog(packageName, "error db set", err)
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
		}
		w.Write([]byte(http.StatusText(http.StatusOK)))
		w.WriteHeader(http.StatusOK)
	}
}

func Wait(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)
	w.WriteHeader(200)
	w.Write([]byte("this is the response"))
}
