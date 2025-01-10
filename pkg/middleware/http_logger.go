package middleware

import (
	"context"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/rlawnsxo131/ws-placeholder/pkg"
	"github.com/rs/zerolog"
)

var (
	HTTPLogEntryCtxKey     = &contextKey{"HTTPLogEntryCtxKey"}
	DefaultHTTPServeLogger = NewHTTPServeLogger(os.Stdout, NewDefaultHTTPLogFormatter())
)

func HTTPLogger(logger *HTTPServeLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			entry := logger.NewLogEntry(r)
			writer := NewHTTPLogResponseWriter(w, entry)

			t := time.Now()
			defer func() {
				entry.Write(t)
			}()

			next.ServeHTTP(writer, WithHTTPLogEntry(r, entry))
		})
	}
}

type HTTPServeLogger struct {
	l *zerolog.Logger
	f HTTPLogFormatter
}

func NewHTTPServeLogger(w io.Writer, f HTTPLogFormatter) *HTTPServeLogger {
	l := zerolog.New(w).With().Caller().Logger()
	return &HTTPServeLogger{
		l: &l,
		f: f,
	}
}

func (l *HTTPServeLogger) NewLogEntry(r *http.Request) HTTPLogEntry {
	return l.f.NewLogEntry(l.l, r)
}

type HTTPLogFormatter interface {
	NewLogEntry(l *zerolog.Logger, r *http.Request) HTTPLogEntry
}

type HTTPLogEntry interface {
	Add(f func(e *zerolog.Event))
	Write(t time.Time)
}

func GetHTTPLogEntry(ctx context.Context) HTTPLogEntry {
	entry, _ := ctx.Value(HTTPLogEntryCtxKey).(HTTPLogEntry)
	return entry
}

func WithHTTPLogEntry(r *http.Request, entry HTTPLogEntry) *http.Request {
	r = r.WithContext(context.WithValue(r.Context(), HTTPLogEntryCtxKey, entry))
	return r
}

type DefaultHTTPLogFormatter struct {
	HTTPLogFormatter
}

func NewDefaultHTTPLogFormatter() *DefaultHTTPLogFormatter {
	return &DefaultHTTPLogFormatter{}
}

func (f *DefaultHTTPLogFormatter) NewLogEntry(l *zerolog.Logger, r *http.Request) HTTPLogEntry {
	return &DefaultHTTPLogEntry{
		l:   l,
		r:   r,
		add: []func(e *zerolog.Event){},
	}
}

type DefaultHTTPLogEntry struct {
	l   *zerolog.Logger
	r   *http.Request
	add []func(e *zerolog.Event)
}

func (le *DefaultHTTPLogEntry) Add(f func(e *zerolog.Event)) {
	le.add = append(le.add, f)
}

func (le *DefaultHTTPLogEntry) Write(t time.Time) {
	e := le.l.Log().
		Str("time", t.UTC().Format(time.RFC3339Nano)).
		Str("request-id", GetHTTPRequestID(le.r.Context())).
		Dur("elapsed(ms)", time.Since(t)).
		Str("method", le.r.Method).
		Str("uri", le.r.RequestURI).
		Str("origin", le.r.Header.Get(pkg.HeaderOrigin)).
		Str("host", le.r.Host).
		Str("referer", le.r.Referer()).
		Str("remote-ip", le.r.RemoteAddr).
		Str("x-request-id", le.r.Header.Get(pkg.HeaderXRequestID)).
		Str("x-forwarded-for", le.r.Header.Get(pkg.HeaderXForwardedFor)).
		Str("cookie", le.r.Header.Get(pkg.HeaderCookie))

	for _, f := range le.add {
		f(e)
	}

	e.Send()
}

type HTTPLogResponseWriter struct {
	w  http.ResponseWriter
	le HTTPLogEntry
}

func NewHTTPLogResponseWriter(w http.ResponseWriter, le HTTPLogEntry) http.ResponseWriter {
	return &HTTPLogResponseWriter{
		w:  w,
		le: le,
	}
}

func (lw *HTTPLogResponseWriter) Write(buf []byte) (int, error) {
	lw.le.Add(func(e *zerolog.Event) {
		e.Bytes("response", buf)
	})
	return lw.w.Write(buf)
}

func (lw *HTTPLogResponseWriter) Header() http.Header {
	return lw.w.Header()
}

func (lw *HTTPLogResponseWriter) WriteHeader(statusCode int) {
	lw.le.Add(func(e *zerolog.Event) {
		e.Int("status", statusCode)
	})
	lw.w.WriteHeader(statusCode)
}
