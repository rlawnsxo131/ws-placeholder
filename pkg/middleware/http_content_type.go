package middleware

import (
	"net/http"

	"github.com/rlawnsxo131/ws-placeholder/pkg/constants"
)

func HTTPContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if contentType := r.Header.Get(constants.HeaderContentType); contentType != "" {
			w.Header().Set(constants.HeaderContentType, contentType)
		} else {
			w.Header().Set(constants.HeaderContentType, "application/json; charset=utf-8")
		}
		next.ServeHTTP(w, r)
	})
}
