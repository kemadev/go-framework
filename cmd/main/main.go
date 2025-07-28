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

	r.Handle("GET /foo", http.NotFoundHandler())

	server.Run(r)
}
