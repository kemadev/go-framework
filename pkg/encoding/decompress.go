// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package encoding

import (
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"sync"

	"github.com/kemadev/go-framework/pkg/convenience/headkey"
	"github.com/kemadev/go-framework/pkg/convenience/headval"
	"github.com/kemadev/go-framework/pkg/convenience/log"
)

// DecompressMiddleware returns a middleware that performs automatic decompression of request body when
// content encoding is gzip. It is inspired from [echo's implementation].
//
// [echo's implementation]: https://github.com/labstack/echo/blob/master/middleware/decompress.go
func DecompressMiddleware(next http.Handler) http.Handler {
	decompressPool := gzipDecompressPool()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(headkey.ContentEncoding) != headval.EncodingGzip {
			next.ServeHTTP(w, r)

			return
		}

		pe := decompressPool.Get()
		defer decompressPool.Put(pe)

		decompressReader, ok := pe.(*gzip.Reader)
		if !ok || decompressReader == nil {
			log.ErrLog(packageName, "error getting decompressor", ErrFailureGetFromPool)
			http.Error(
				w,
				http.StatusText(http.StatusServiceUnavailable),
				http.StatusServiceUnavailable,
			)

			return
		}

		body := r.Body
		defer body.Close()

		err := decompressReader.Reset(body)
		if err != nil {
			// Ignore empty body errors
			if errors.Is(err, io.EOF) {
				next.ServeHTTP(w, r)

				return
			}

			var maxBytesErr *http.MaxBytesError
			if errors.As(err, &maxBytesErr) {
				http.Error(
					w,
					http.StatusText(http.StatusRequestEntityTooLarge),
					http.StatusRequestEntityTooLarge,
				)

				return
			}

			log.ErrLog(packageName, "error resetting body decompressor: %w", err)
			http.Error(
				w,
				http.StatusText(http.StatusServiceUnavailable),
				http.StatusServiceUnavailable,
			)

			return
		}

		// Only Close gzip reader if it was set to a proper gzip source to prevent panic on closure
		defer decompressReader.Close()

		r.Body = decompressReader

		next.ServeHTTP(w, r)
	})
}

func gzipDecompressPool() sync.Pool {
	return sync.Pool{
		New: func() any {
			r := new(gzip.Reader)

			return r
		},
	}
}
