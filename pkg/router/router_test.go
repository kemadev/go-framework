// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChain(t *testing.T) {
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

	handler := func(w http.ResponseWriter, r *http.Request) {}

	mux := http.NewServeMux()

	chain1 := chain{mw1, mw2}

	mux.Handle("GET /{$}", chain1.thenFunc(handler))

	chain2 := append(chain1, mw3, mw4)
	mux.Handle("GET /foo", chain2.thenFunc(handler))

	chain3 := append(chain2, mw5)
	mux.Handle("GET /nested/foo", chain3.thenFunc(handler))

	chain4 := append(chain1, mw6)
	mux.Handle("GET /bar", chain4.thenFunc(handler))

	mux.Handle("GET /baz", chain1.thenFunc(handler))

	tests := []struct {
		RequestMethod  string
		RequestPath    string
		ExpectedUsed   string
		ExpectedStatus int
	}{
		{
			RequestMethod:  "GET",
			RequestPath:    "/",
			ExpectedUsed:   "12",
			ExpectedStatus: http.StatusOK,
		},
		{
			RequestMethod:  "GET",
			RequestPath:    "/foo",
			ExpectedUsed:   "1234",
			ExpectedStatus: http.StatusOK,
		},
		{
			RequestMethod:  "GET",
			RequestPath:    "/nested/foo",
			ExpectedUsed:   "12345",
			ExpectedStatus: http.StatusOK,
		},
		{
			RequestMethod:  "GET",
			RequestPath:    "/bar",
			ExpectedUsed:   "126",
			ExpectedStatus: http.StatusOK,
		},
		{
			RequestMethod:  "GET",
			RequestPath:    "/baz",
			ExpectedUsed:   "12",
			ExpectedStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		used = ""

		req, err := http.NewRequest(test.RequestMethod, test.RequestPath, nil)
		if err != nil {
			t.Errorf("NewRequest: %s", err)
		}

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

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
