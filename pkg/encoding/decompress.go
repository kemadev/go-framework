package encoding

import (
	"compress/gzip"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"sync"

	"github.com/kemadev/go-framework/pkg/convenience/headkey"
	"github.com/kemadev/go-framework/pkg/convenience/headval"
	"github.com/kemadev/go-framework/pkg/convenience/log"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

var ErrFailureGetDecompressorFromPool = errors.New("can't get a decompressor from pool")

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
			log.Logger(packageName).
				Error("error decompressing body", slog.String(string(semconv.ErrorMessageKey), ErrFailureGetDecompressorFromPool.Error()))
			return
		}

		body := r.Body
		defer body.Close()

		err := decompressReader.Reset(body)
		if err != nil {
			// Ignore empty body errors
			if err == io.EOF {
				next.ServeHTTP(w, r)
				return
			}
			log.Logger(packageName).
				Error("error resetting body decompressor: %w", slog.String(string(semconv.ErrorMessageKey), err.Error()))
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
		New: func() interface{} {
			r := new(gzip.Reader)
			return r
		},
	}
}
