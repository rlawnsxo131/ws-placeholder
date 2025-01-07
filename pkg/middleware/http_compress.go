package middleware

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/rlawnsxo131/ws-placeholder/pkg"
	"github.com/rlawnsxo131/ws-placeholder/pkg/lib/logger"
)

func HTTPCompress(level int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		gzipScheme := "gzip"
		gzipPool := gzipCompressPool(level)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("compress")
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

			next.ServeHTTP(&HTTPCompressWriter{Writer: gw, ResponseWriter: w}, r)
		})
	}
}

type HTTPCompressWriter struct {
	io.Writer
	http.ResponseWriter
}

func (cw *HTTPCompressWriter) Write(buf []byte) (int, error) {
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
