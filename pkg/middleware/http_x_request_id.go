package middleware

import (
	"strings"

	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rlawnsxo131/ws-placeholder/pkg/constants"

	"net/http"
)

func HTTPXRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// chi_middleware.GetReqID(r.Context())
		// ex) hostname/qzmQUuE5WX-000001
		if id := chi_middleware.GetReqID(r.Context()); id != "" {
			splitId := strings.Split(id, "/")[1]
			r.Header.Set(constants.HeaderXRequestID, splitId)
		}
		next.ServeHTTP(w, r)
	})
}
