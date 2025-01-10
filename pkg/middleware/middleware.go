package middleware

// [Full order example]
// r.Use(middleware.HTTPRecovery) // defer 3
// r.Use(middleware.HTTPCompress(middleware.HTTPCompressConfig{ ... })) // defer 2: clean up gzipWriter
// r.Use(middleware.HTTPRequestID) // context set id
// r.Use(middleware.HTTPLogger(middleware.DefaultHTTPServeLogger)) // defer 1
// r.Use(middleware.HTTPRealIP)
// r.Use(middleware.HTTPTimeout(time.Second * 3))
// r.Use(middleware.HTTPCors(middleware.HTTPCorsConfig{ ... }))
// r.Use(middleware.HTTPContentType(middleware.HeaderJson))

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "pkg/middleware context value " + k.name
}
