package middleware

// [Full order example]
// r.Use(middleware.HTTPCompress(5))
// r.Use(middleware.HTTPLogger(middleware.DefaultHTTPServeLogger))
// r.Use(middleware.HTTPRequestID)
// r.Use(middleware.HTTPRealIP)
// r.Use(middleware.HTTPTimeout(time.Second * 3))
// r.Use(middleware.HTTPCors(middleware.HTTPCorsConfig{ ... }))
// r.Use(middleware.HTTPContentType(middleware.HeaderJson))
// r.Use(middleware.HTTPRecovery)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "pkg/middleware context value " + k.name
}
