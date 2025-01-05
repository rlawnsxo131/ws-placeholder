package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/rs/zerolog"
)

func HTTPRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if err == http.ErrAbortHandler {
					// we don't recover http.ErrAbortHandler so the response
					// to the client is aborted, this should not be logged
					panic(err)
				}
				GetHTTPLogEntry(r).Add(func(e *zerolog.Event) {
					e.Any("recover panic err", err).Str("statck", string(debug.Stack()))
				})
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			}
		}()
		next.ServeHTTP(w, r)
	})
}