package middleware

import (
	"net/http"

	"github.com/rlawnsxo131/ws-placeholder/pkg"
)

type HeaderContentType int

const (
	HeaderJson HeaderContentType = iota + 1
	HeaderText
)

var contentTypes = []string{
	"application/json; charset=utf-8",
	"plain/text; charset=utf-8",
}

func (contentType HeaderContentType) String() string {
	return contentTypes[contentType-1]
}

// https://www.iana.org/assignments/media-types/media-types.xhtml
func HTTPContentType(contentType HeaderContentType) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(pkg.HeaderContentType, HeaderContentType(contentType).String())
			next.ServeHTTP(w, r)
		})
	}
}
