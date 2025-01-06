package middleware

import (
	"net"
	"net/http"
	"strings"

	"github.com/rlawnsxo131/ws-placeholder/pkg/constants"
)

func HTTPRealIP(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if rip := realIP(r); rip != "" {
			r.RemoteAddr = rip
		}
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

var trueClientIP = http.CanonicalHeaderKey(constants.HeaderTrueClientIP)
var xForwardedFor = http.CanonicalHeaderKey(constants.HeaderXForwardedFor)
var xRealIP = http.CanonicalHeaderKey(constants.HeaderXRealIP)

func realIP(r *http.Request) string {
	var ip string

	if tcip := r.Header.Get(trueClientIP); tcip != "" {
		ip = tcip
	} else if xrip := r.Header.Get(xRealIP); xrip != "" {
		ip = xrip
	} else if xff := r.Header.Get(xForwardedFor); xff != "" {
		i := strings.Index(xff, ",")
		if i == -1 {
			i = len(xff)
		}
		ip = xff[:i]
	}
	if ip == "" || net.ParseIP(ip) == nil {
		return ""
	}
	return ip
}
