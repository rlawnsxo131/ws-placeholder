package middleware

// [Full order example]
// r.Use(middleware.HTTPRequestID)
// r.Use(middleware.HTTPXRequestID)
// r.Use(middleware.HTTPRealIP)
// r.Use(middleware.HTTPCompress(5))
// r.Use(middleware.HTTPLogger(middleware.DefaultHTTPServeLogger))
// r.Use(middleware.HTTPTimeout(time.Second * 3))
// r.Use(middleware.HTTPContentType)
// r.Use(middleware.HTTPRecoverer)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "pkg/middleware context value " + k.name
}
