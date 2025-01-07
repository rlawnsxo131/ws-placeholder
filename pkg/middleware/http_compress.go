package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/rlawnsxo131/ws-placeholder/pkg"
	"github.com/rlawnsxo131/ws-placeholder/pkg/lib/logger"
)

// @TODO 다시 다듬기
func HTTPCompress(cfg HTTPCompressConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		gzipScheme := "gzip"
		gzipPool := gzipCompressPool(cfg.Level)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(pkg.HeaderVary, pkg.HeaderAcceptEncoding)

			acceptEncoding := r.Header.Get(pkg.HeaderAcceptEncoding)

			if !strings.Contains(acceptEncoding, gzipScheme) {
				next.ServeHTTP(w, r)
				return
			}

			w.Header().Set(pkg.HeaderContentEncoding, gzipScheme)

			gw, ok := gzipPool.Get().(*gzip.Writer)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}
			gw.Reset(w)

			defer func() {
				if err := gw.Close(); err != nil {
					logger.Default().Err(err).Send()
				}
				gzipPool.Put(gw)
			}()

			next.ServeHTTP(
				&HTTPCompressWriter{
					Writer:         gw,
					ResponseWriter: w,
					minLength:      cfg.MinLength,
				},
				r,
			)
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
	minLength int
}

func (cw *HTTPCompressWriter) Write(buf []byte) (int, error) {
	if cw.Header().Get(pkg.HeaderContentType) == "" {
		cw.Header().Set(pkg.HeaderContentType, http.DetectContentType(buf))
	}
	return cw.Writer.Write(buf)
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

func bufferPool() sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			b := &bytes.Buffer{}
			return b
		},
	}
}
