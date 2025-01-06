package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/rlawnsxo131/ws-placeholder/pkg/constants"
	"github.com/rlawnsxo131/ws-placeholder/pkg/lib/logger"
)

const (
	_gzipScheme    = "gzip"
	_deflateScheme = "deflate"
)

// @TODO config 구현
// deflate 지원할까 말까
func HTTPCompress(level int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		gzipPool := gzipCompressPool(level)
		// deflatePool := deflateCompressPool(level)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(constants.HeaderVary, constants.HeaderAcceptEncoding)

			acceptEncoding := r.Header.Get(constants.HeaderAcceptEncoding)

			if !strings.Contains(acceptEncoding, _gzipScheme) {
				next.ServeHTTP(w, r)
				return
			}

			w.Header().Set(constants.HeaderContentEncoding, _gzipScheme)

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

// func deflateCompressPool(level int) sync.Pool {
// 	return sync.Pool{
// 		New: func() any {
// 			w, err := flate.NewWriter(io.Discard, level)
// 			if err != nil {
// 				return err
// 			}
// 			return w
// 		},
// 	}
// }
