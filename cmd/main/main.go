/*
Copyright 2025 kemadev
SPDX-License-Identifier: MPL-2.0
*/

package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kemadev/go-framework/pkg/convenience/otel"
	"github.com/kemadev/go-framework/pkg/server"
)

const packageName = "github.com/kemadev/go-framework/cmd/main"

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get(otel.WrapHandler("/articles/{rid:^[0-9]{5,6}}", basicHandler))

	server.Run(otel.WrapMux(r, packageName))
}

func basicHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(
		[]byte(
			"Pattern: " + chi.URLParam(
				r,
				"rid",
			) + " - frompath=" + r.PathValue(
				"rid",
			) + " q=" + r.URL.Query().
				Get("q"),
		),
	)
}
