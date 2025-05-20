package logging

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/felixge/httpsnoop"
)

type SlogMiddleware struct {
	log *slog.Logger
}

type SlogOptionFunc func(*SlogMiddleware)

func WithSlog(l *slog.Logger) func(*SlogMiddleware) {
	return func(sm *SlogMiddleware) {
		sm.log = l
	}
}

func NewSlog(opts ...SlogOptionFunc) *SlogMiddleware {
	l := &SlogMiddleware{
		log: slog.Default(),
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func (m *SlogMiddleware) Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		// Execute the chain of handlers, while capturing HTTP metrics: code, bytes-written, duration
		metrics := httpsnoop.CaptureMetrics(next, w, r)

		host := r.Header.Get("x-forwarded-for")
		if host == "" {
			// r.RemoteAddr contains port, which we want to remove
			idx := strings.LastIndex(r.RemoteAddr, ":")
			if idx == -1 {
				host = r.RemoteAddr
			} else {
				host = r.RemoteAddr[:idx]
			}
		}

		m.log.Info("",
			"host", host,
			"method", r.Method,
			"path", r.URL.Path,
			"proto", r.Proto,
			"status", metrics.Code,
			"bytes", metrics.Written,
			"referer", r.Header.Get("referer"),
			"user_agent", r.Header.Get("user-agent"),
			"time_ms", metrics.Duration.Milliseconds(),
		)
	}

	return http.HandlerFunc(fn)
}
