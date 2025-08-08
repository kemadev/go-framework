package encoding

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"sync"

	"github.com/kemadev/go-framework/pkg/convenience/headkey"
	"github.com/kemadev/go-framework/pkg/convenience/headutil"
	"github.com/kemadev/go-framework/pkg/convenience/headval"
	"github.com/kemadev/go-framework/pkg/convenience/log"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

// CompressConfig defines the configuration for compression middleware.
type CompressConfig struct {
	// Gzip compression level (-1 to 9, where -1 is default compression)
	Level int

	// Minimum response length before compression is applied
	MinLength int
}

type compressResponseWriter struct {
	http.ResponseWriter
	writer            *gzip.Writer
	buffer            *bytes.Buffer
	minLength         int
	wroteHeader       bool
	wroteBody         bool
	minLengthExceeded bool
	statusCode        int
}

// gzip overhead makes response bigger for small bodies, and thus wastes CPU time for counter-productive results
const CompressionMinThreshold = 2 * 1024

// CompressMiddleware returns a middleware that performs automatic compression of response body
// when the client accepts gzip encoding.
func CompressMiddleware(next http.Handler) http.Handler {
	return CompressMiddlewareWithConfig(CompressConfig{
		Level: gzip.DefaultCompression,
		MinLength: CompressionMinThreshold,
	})(next)
}

// CompressMiddlewareWithConfig returns a middleware with custom compression configuration.
func CompressMiddlewareWithConfig(conf CompressConfig) func(http.Handler) http.Handler {
	if conf.Level == 0 {
		conf.Level = gzip.DefaultCompression
	}

	compressPool := gzipCompressPool(conf)
	bufferPool := bufferPool()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(headkey.Vary, headkey.AcceptEncoding)

			if !headutil.AcceptsEncoding(r.Header, headval.EncodingGzip) {
				next.ServeHTTP(w, r)
				return
			}

			pe := compressPool.Get()
			defer compressPool.Put(pe)

			gzipWriter, ok := pe.(*gzip.Writer)
			if !ok || gzipWriter == nil {
				log.Logger(packageName).
					Error("error getting compressor from pool",
						slog.String(string(semconv.ErrorMessageKey), "invalid compressor type"))
				next.ServeHTTP(w, r)

				return
			}

			gzipWriter.Reset(w)

			be := bufferPool.Get()
			defer bufferPool.Put(be)

			buffer, ok := be.(*bytes.Buffer)
			if !ok || buffer == nil {
				log.Logger(packageName).
					Error("error getting buffer from pool",
						slog.String(string(semconv.ErrorMessageKey), "invalid buffer type"))
				next.ServeHTTP(w, r)

				return
			}

			buffer.Reset()

			crw := &compressResponseWriter{
				ResponseWriter:    w,
				writer:            gzipWriter,
				buffer:            buffer,
				minLength:         conf.MinLength,
				wroteHeader:       false,
				wroteBody:         false,
				minLengthExceeded: false,
				statusCode:        http.StatusOK,
			}

			defer func() {
				if !crw.wroteBody {
					if crw.wroteHeader {
						crw.ResponseWriter.WriteHeader(crw.statusCode)
					}
				} else if !crw.minLengthExceeded {
					if crw.wroteHeader {
						crw.ResponseWriter.WriteHeader(crw.statusCode)
					}
					_, err := crw.buffer.WriteTo(crw.ResponseWriter)
					if err != nil {
						log.Logger(packageName).
							Error("error writing uncompressed response",
								slog.String(string(semconv.ErrorMessageKey), err.Error()))
					}
					gzipWriter.Reset(io.Discard)
				}

				err := gzipWriter.Close()
				if err != nil {
					log.Logger(packageName).
						Error("error closing gzip writer",
							slog.String(string(semconv.ErrorMessageKey), err.Error()))
				}

				compressPool.Put(gzipWriter)
				bufferPool.Put(buffer)
			}()

			next.ServeHTTP(crw, r)
		})
	}
}

// WriteHeader implements [net/http.ResponseWriter].
func (w *compressResponseWriter) WriteHeader(statusCode int) {
	// Remove Content-Length header as it will be invalid after compression
	w.Header().Del(headkey.ContentLength)

	w.wroteHeader = true
	w.statusCode = statusCode
}

// Write implements [net/http.ResponseWriter].
func (w *compressResponseWriter) Write(data []byte) (int, error) {
	w.wroteBody = true

	// If not exceeding minimum length, buffer
	if !w.minLengthExceeded {
		n, err := w.buffer.Write(data)
		if err != nil {
			return n, err
		}

		if w.buffer.Len() >= w.minLength {
			w.minLengthExceeded = true

			w.Header().Set(headkey.ContentEncoding, headval.EncodingGzip)
			if w.wroteHeader {
				w.ResponseWriter.WriteHeader(w.statusCode)
			}

			return w.writer.Write(w.buffer.Bytes())
		}

		return n, nil
	}

	// We already decided to compress, just write
	n, err := w.writer.Write(data)
	if err != nil {
		return n, fmt.Errorf("error writing compressed response: %w", err)
	}

	return n, nil
}

func (w *compressResponseWriter) Flush() {
	// If we haven't exceeded minimum length, force compression
	if !w.minLengthExceeded {
		w.minLengthExceeded = true

		w.Header().Set(headkey.ContentEncoding, headval.EncodingGzip)
		if w.wroteHeader {
			w.ResponseWriter.WriteHeader(w.statusCode)
		}

		_, err := w.Write(w.buffer.Bytes())
		if err != nil {
			log.Logger(packageName).
				Error("error writing buffered data during flush",
					slog.String(string(semconv.ErrorMessageKey), err.Error()))
		}
	}

	err := w.writer.Flush()
	if err != nil {
		log.Logger(packageName).
			Error("error flushing gzip writer",
				slog.String(string(semconv.ErrorMessageKey), err.Error()))
	}

	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Hijack implements [net/http.Hijacker] if the underlying ResponseWriter supports it.
func (w *compressResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}

	return hijacker.Hijack()
}

// Push implements [net/http.Pusher] if the underlying ResponseWriter supports it.
func (w *compressResponseWriter) Push(target string, opts *http.PushOptions) error {
	pusher, ok := w.ResponseWriter.(http.Pusher)
	if !ok {
		return http.ErrNotSupported
	}

	return pusher.Push(target, opts)
}

// Unwrap returns the underlying [net/http.ResponseWriter].
func (w *compressResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func gzipCompressPool(conf CompressConfig) sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			w, err := gzip.NewWriterLevel(io.Discard, conf.Level)
			if err != nil {
				return err
			}
			return w
		},
	}
}

func bufferPool() sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			b := &bytes.Buffer{}
			return b
		},
	}
}
