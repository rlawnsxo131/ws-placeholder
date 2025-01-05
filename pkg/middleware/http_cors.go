// was modified as needed by referring to cors middleware at https://github.com/labstack/echo

package middleware

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/rlawnsxo131/ws-placeholder/pkg/constants"
	"github.com/rlawnsxo131/ws-placeholder/pkg/lib/logger"
)

var DefaultHTTPCorsConfig = HTTPCorsConfig{
	AllowOrigins: []string{"*"},
	AllowMethods: []string{http.MethodGet},
	AllowHeaders: []string{
		constants.HeaderContentType,
		constants.HeaderAccept,
	},
}

func HTTPCors(config HTTPCorsConfig) func(http.Handler) http.Handler {
	allowOriginPatterns := make([]*regexp.Regexp, 0, len(config.AllowOrigins))
	for _, origin := range config.AllowOrigins {
		if origin == "*" {
			continue // "*" is handled differently and does not need regexp
		}
		pattern := regexp.QuoteMeta(origin)
		pattern = strings.ReplaceAll(pattern, "\\*", ".*")
		pattern = strings.ReplaceAll(pattern, "\\?", ".")
		pattern = "^" + pattern + "$"

		re, err := regexp.Compile(pattern)
		if err != nil {
			logger.Default().Panic().Err(err).Send()
		}
		allowOriginPatterns = append(allowOriginPatterns, re)
	}

	allowMethods := strings.Join(config.AllowMethods, ",")
	allowHeaders := strings.Join(config.AllowHeaders, ",")
	exposeHeaders := strings.Join(config.ExposeHeaders, ",")

	maxAge := "0"
	if config.MaxAge > 0 {
		maxAge = strconv.Itoa(config.MaxAge)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get(constants.HeaderOrigin)
			allowOrigin := ""

			w.Header().Add(constants.HeaderVary, constants.HeaderOrigin)

			// Preflight request is an OPTIONS request, using three HTTP request headers: Access-Control-Request-Method,
			// Access-Control-Request-Headers, and the Origin header. See: https://developer.mozilla.org/en-US/docs/Glossary/Preflight_request
			// For simplicity we just consider method type and later `Origin` header.
			preflight := r.Method == http.MethodOptions

			if origin == "" {
				if !preflight {
					next.ServeHTTP(w, r)
					return
				}
				w.WriteHeader(http.StatusNoContent)
				w.Write([]byte(http.StatusText(http.StatusNoContent)))
				return
			}
			if !preflight && r.Header.Get(constants.HeaderContentType) == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(http.StatusText(http.StatusBadRequest)))
				return
			}

			for _, o := range config.AllowOrigins {
				if o == "*" && config.AllowCredentials {
					allowOrigin = origin
					break
				}
				if o == "*" || o == origin {
					allowOrigin = o
					break
				}
				if matchSubdomain(origin, o) {
					allowOrigin = origin
					break
				}
			}

			checkPatterns := false
			if allowOrigin == "" {
				// to avoid regex cost by invalid (long) domains (253 is domain name max limit)
				if len(origin) <= (253+3+5) && strings.Contains(origin, "://") {
					checkPatterns = true
				}
			}
			if checkPatterns {
				for _, re := range allowOriginPatterns {
					if match := re.MatchString(origin); match {
						allowOrigin = origin
						break
					}
				}
			}

			// Origin not allowed
			if allowOrigin == "" {
				if !preflight {
					next.ServeHTTP(w, r)
					return
				}
				w.WriteHeader(http.StatusNoContent)
				w.Write([]byte(http.StatusText(http.StatusNoContent)))
				return
			}

			w.Header().Set(constants.HeaderAccessControlAllowOrigin, allowOrigin)
			if config.AllowCredentials {
				w.Header().Set(constants.HeaderAccessControlAllowCredentials, "true")
			}

			// Simple request
			if !preflight {
				if exposeHeaders != "" {
					w.Header().Set(constants.HeaderAccessControlExposeHeaders, exposeHeaders)
				}
				next.ServeHTTP(w, r)
				return
			}

			// Preflight request
			w.Header().Add(constants.HeaderVary, constants.HeaderAccessControlRequestMethod)
			w.Header().Add(constants.HeaderVary, constants.HeaderAccessControlRequestHeaders)
			w.Header().Set(constants.HeaderAccessControlAllowMethods, allowMethods)

			if allowHeaders != "" {
				w.Header().Set(constants.HeaderAccessControlAllowHeaders, allowHeaders)
			} else {
				h := r.Header.Get(constants.HeaderAccessControlRequestHeaders)
				if h != "" {
					w.Header().Set(constants.HeaderAccessControlAllowHeaders, h)
				}
			}
			if config.MaxAge != 0 {
				w.Header().Set(constants.HeaderAccessControlMaxAge, maxAge)
			}

			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte(http.StatusText(http.StatusNoContent)))
		})
	}

}

type HTTPCorsConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	ExposeHeaders    []string
	MaxAge           int
}

func matchScheme(domain, pattern string) bool {
	didx := strings.Index(domain, ":")
	pidx := strings.Index(pattern, ":")
	return didx != -1 && pidx != -1 && domain[:didx] == pattern[:pidx]
}

// matchSubdomain compares authority with wildcard
func matchSubdomain(domain, pattern string) bool {
	if !matchScheme(domain, pattern) {
		return false
	}
	didx := strings.Index(domain, "://")
	pidx := strings.Index(pattern, "://")
	if didx == -1 || pidx == -1 {
		return false
	}
	domAuth := domain[didx+3:]
	// to avoid long loop by invalid long domain
	if len(domAuth) > 253 {
		return false
	}
	patAuth := pattern[pidx+3:]

	domComp := strings.Split(domAuth, ".")
	patComp := strings.Split(patAuth, ".")
	for i := len(domComp)/2 - 1; i >= 0; i-- {
		opp := len(domComp) - 1 - i
		domComp[i], domComp[opp] = domComp[opp], domComp[i]
	}
	for i := len(patComp)/2 - 1; i >= 0; i-- {
		opp := len(patComp) - 1 - i
		patComp[i], patComp[opp] = patComp[opp], patComp[i]
	}

	for i, v := range domComp {
		if len(patComp) <= i {
			return false
		}
		p := patComp[i]
		if p == "*" {
			return true
		}
		if p != v {
			return false
		}
	}
	return false
}
