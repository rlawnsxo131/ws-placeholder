package middleware

import (
	"net/http"
	"time"
)

func HTTPTimeout(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler := http.TimeoutHandler(next, timeout, http.StatusText(http.StatusGatewayTimeout))
			handler.ServeHTTP(w, r)
		})
	}
}
