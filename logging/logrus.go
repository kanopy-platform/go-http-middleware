package logging

import (
	"net/http"
	"strings"

	"github.com/felixge/httpsnoop"
	log "github.com/sirupsen/logrus"
)

type LogrusMiddleware struct {
	log *log.Logger
}

type LogrusOptionFunc func(*LogrusMiddleware)

func WithLogrus(l *log.Logger) func(*LogrusMiddleware) {
	return func(lm *LogrusMiddleware) {
		lm.log = l
	}
}

func NewLogrus(opts ...LogrusOptionFunc) *LogrusMiddleware {
	l := &LogrusMiddleware{
		log: log.StandardLogger(),
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func (m *LogrusMiddleware) Middleware(next http.Handler) http.Handler {
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

		m.log.WithFields(log.Fields{
			"host":       host,
			"method":     r.Method,
			"path":       r.URL.Path,
			"proto":      r.Proto,
			"status":     metrics.Code,
			"bytes":      metrics.Written,
			"referer":    r.Header.Get("referer"),
			"user_agent": r.Header.Get("user-agent"),
			"time_ms":    metrics.Duration.Milliseconds(),
		}).Info("handled")
	}
	return http.HandlerFunc(fn)
}
