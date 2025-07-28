/*
Copyright 2025 kemadev
SPDX-License-Identifier: MPL-2.0
*/

package main

import (
	"net/http"

	"github.com/kemadev/go-framework/pkg/router"
	"github.com/kemadev/go-framework/pkg/server"
)

func main() {
	r := router.New()

	// Default handler when nothing matches
	r.HandleOTEL("/", http.NotFoundHandler())

	r.HandleOTEL("GET /foo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	}))

	server.Run(r)
}
