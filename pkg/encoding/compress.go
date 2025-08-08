package encoding

import (
	"compress/gzip"
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"github.com/kemadev/go-framework/pkg/convenience/headutil"
	"github.com/kemadev/go-framework/pkg/convenience/headval"
	"github.com/kemadev/go-framework/pkg/convenience/log"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

var ErrFailureGetCompressorFromPool = errors.New("can't get a compressor from pool")

type Compressor interface {
	gzipCompressPool() sync.Pool
}

var CompressPool = sync.Pool{
	New: func() any {
		return new(gzip.Writer)
	},
}

// DecompressMiddleware returns a middleware that performs automatic compression of response body
// using gzip. It is inspired from [echo's implementation].
//
// [echo's implementation]: https://github.com/labstack/echo/blob/master/middleware/compress.go
func CompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !headutil.AcceptsEncoding(r.Header, headval.EncodingGzip) {
			next.ServeHTTP(w, r)
			return
		}

		pe := CompressPool.Get()
		defer CompressPool.Put(pe)

		compressWriter, ok := pe.(*gzip.Writer)
		if !ok || compressWriter == nil {
			log.Logger(packageName).
				Error("error compressing body", slog.String(string(semconv.ErrorMessageKey), ErrFailureGetCompressorFromPool.Error()))
			return
		}

		defer compressWriter.Close()

		compressWriter.Reset(w)

		next.ServeHTTP(w, r)
	})
}

func compressionDesirable(mime string) bool {
	// Ceretains mimes should compress only, not all
	// Small bodies should not compress
	return false
}
