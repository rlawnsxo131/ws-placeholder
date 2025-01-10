package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/rlawnsxo131/ws-placeholder/pkg"
	"github.com/rlawnsxo131/ws-placeholder/pkg/lib/logger"
)

const _gzipScheme = "gzip"

func HTTPCompress(cfg HTTPCompressConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		gzipPool := gzipCompressPool(cfg.Level)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(pkg.HeaderVary, pkg.HeaderAcceptEncoding)

			acceptEncoding := r.Header.Get(pkg.HeaderAcceptEncoding)

			if !strings.Contains(acceptEncoding, _gzipScheme) {
				next.ServeHTTP(w, r)
				return
			}

			gw, ok := gzipPool.Get().(*gzip.Writer)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}
			gw.Reset(w)

			cw := &HTTPCompressWriter{
				Writer:         gw,
				ResponseWriter: w,
				minLength:      cfg.MinLength,
			}

			defer func() {
				if cw.compressed {
					if err := gw.Close(); err != nil {
						logger.Default().Err(err).Send()
					}
				}
				gzipPool.Put(gw)
			}()

			next.ServeHTTP(cw, r)
		})
	}
}

type HTTPCompressConfig struct {
	Level     int
	MinLength int
}

type HTTPCompressWriter struct {
	io.Writer
	http.ResponseWriter
	minLength   int
	statusCode  int
	wroteHeader bool
	compressed  bool
}

func (cw *HTTPCompressWriter) Write(buf []byte) (int, error) {
	if cw.minLength <= len(buf) {
		cw.compressed = true
		if cw.wroteHeader {
			cw.WriteHeader(cw.statusCode)
		}
		cw.Header().Set(pkg.HeaderContentEncoding, _gzipScheme)
		return cw.Writer.Write(buf)
	}

	// uncompressed
	cw.compressed = false
	if cw.wroteHeader {
		cw.ResponseWriter.WriteHeader(cw.statusCode)
	}
	if cw.Header().Get(pkg.HeaderContentType) == "" {
		cw.Header().Set(pkg.HeaderContentType, http.DetectContentType(buf))
	}

	return cw.ResponseWriter.Write(buf)
}

func (cw *HTTPCompressWriter) WriteHeader(code int) {
	cw.Header().Del(pkg.HeaderContentLength)

	cw.wroteHeader = true

	// Delay writing of the header until we know if we'll actually compress the response
	cw.statusCode = code
}

func gzipCompressPool(level int) sync.Pool {
	return sync.Pool{
		New: func() any {
			w, err := gzip.NewWriterLevel(io.Discard, level)
			if err != nil {
				return err
			}
			return w
		},
	}
}
