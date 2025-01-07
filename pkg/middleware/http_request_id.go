package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/rlawnsxo131/ws-placeholder/pkg"
)

var HTTPRequestIDKey = &contextKey{"HTTPRequestIDKey"}

func HTTPRequestID(next http.Handler) http.Handler {
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}
	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		rand.Read(buf[:])
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}

	var reqid uint64 = 0
	var prefix string = fmt.Sprintf("%s/%s", hostname, b64[0:10])

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(pkg.HeaderXRequestID)
		if requestID == "" {
			myid := atomic.AddUint64(&reqid, 1)
			requestID = fmt.Sprintf("%s-%06d", prefix, myid)
		}
		// ex) hostname/qzmQUuE5WX-000001
		r.Header.Set(pkg.HeaderXRequestID, strings.Split(requestID, "/")[1])
		next.ServeHTTP(w, WithHTTPRequestID(r, requestID))
	})
}

func GetHTTPRequestID(ctx context.Context) string {
	if id, _ := ctx.Value(HTTPRequestIDKey).(string); id != "" {
		return id
	}
	return ""
}

func WithHTTPRequestID(r *http.Request, id string) *http.Request {
	r = r.WithContext(context.WithValue(r.Context(), HTTPRequestIDKey, id))
	return r
}
