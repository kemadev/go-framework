// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package router_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kemadev/go-framework/pkg/router"
)

func TestChain(t *testing.T) {
	t.Parallel()

	used := ""

	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			used += "1"

			next.ServeHTTP(w, r)
		})
	}

	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			used += "2"

			next.ServeHTTP(w, r)
		})
	}

	mw3 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			used += "3"

			next.ServeHTTP(w, r)
		})
	}

	mw4 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			used += "4"

			next.ServeHTTP(w, r)
		})
	}

	mw5 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			used += "5"

			next.ServeHTTP(w, r)
		})
	}

	mw6 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			used += "6"

			next.ServeHTTP(w, r)
		})
	}

	handler := func(_ http.ResponseWriter, _ *http.Request) {
		used += "h"
	}

	app1 := router.New()

	app1.Use(mw1)
	app1.Use(mw2)

	app1.HandleFunc("GET /{$}", handler)

	app2 := router.New()
	app2.Use(mw1)
	app2.Use(mw2)
	app2.Use(mw3)
	app2.Use(mw4)
	app2.HandleFunc("GET /foo", handler)

	app3 := router.New()
	app3.Use(mw1)
	app3.Use(mw2)
	app3.Use(mw3)
	app3.Use(mw4)
	app3.Use(mw5)
	app3.HandleFunc("GET /nested/foo", handler)

	app4 := router.New()
	app4.Use(mw1)
	app4.Use(mw2)
	app4.Group(func(r *router.Router) {
		r.Use(mw3)
		r.Use(mw4)
		r.Use(mw5)
		r.HandleFunc("GET /bar", handler)
		r.Group(func(sr *router.Router) {
			sr.Use(mw6)
			sr.HandleFunc("GET /qux", handler)
		})
	})
	app4.Group(func(r *router.Router) {
		r.Use(mw6)
		r.HandleFunc("GET /baz", handler)
	})

	tests := []struct {
		Mux            *router.Router
		RequestMethod  string
		RequestPath    string
		ExpectedUsed   string
		ExpectedStatus int
	}{
		{
			Mux:            app1,
			RequestMethod:  "GET",
			RequestPath:    "/",
			ExpectedUsed:   "12h",
			ExpectedStatus: http.StatusOK,
		},
		{
			Mux:            app2,
			RequestMethod:  "GET",
			RequestPath:    "/foo",
			ExpectedUsed:   "1234h",
			ExpectedStatus: http.StatusOK,
		},
		{
			Mux:            app3,
			RequestMethod:  "GET",
			RequestPath:    "/nested/foo",
			ExpectedUsed:   "12345h",
			ExpectedStatus: http.StatusOK,
		},
		{
			Mux:            app4,
			RequestMethod:  "GET",
			RequestPath:    "/bar",
			ExpectedUsed:   "12345h",
			ExpectedStatus: http.StatusOK,
		},
		{
			Mux:            app4,
			RequestMethod:  "GET",
			RequestPath:    "/baz",
			ExpectedUsed:   "126h",
			ExpectedStatus: http.StatusOK,
		},
		{
			Mux:            app4,
			RequestMethod:  "GET",
			RequestPath:    "/qux",
			ExpectedUsed:   "123456h",
			ExpectedStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		used = ""

		req, err := http.NewRequestWithContext(
			context.Background(),
			test.RequestMethod,
			test.RequestPath,
			nil,
		)
		if err != nil {
			t.Errorf("NewRequestWithContext: %s", err)
		}

		rr := httptest.NewRecorder()
		test.Mux.ServeHTTP(rr, req)

		resp := rr.Result()

		if resp.StatusCode != test.ExpectedStatus {
			t.Errorf(
				"%s %s: expected status %d but was %d",
				test.RequestMethod,
				test.RequestPath,
				test.ExpectedStatus,
				resp.StatusCode,
			)
		}

		if used != test.ExpectedUsed {
			t.Errorf(
				"%s %s: middleware used: expected %q; got %q",
				test.RequestMethod,
				test.RequestPath,
				test.ExpectedUsed,
				used,
			)
		}
	}
}
