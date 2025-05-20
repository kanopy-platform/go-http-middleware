package logging

import (
	"maps"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

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
		start := time.Now()
		recorder := httptest.NewRecorder()

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

		next.ServeHTTP(recorder, r)

		maps.Copy(w.Header(), recorder.Header())
		w.WriteHeader(recorder.Code)
		w.Write(recorder.Body.Bytes())

		duration := time.Since(start)

		m.log.WithFields(log.Fields{
			"host":       host,
			"method":     r.Method,
			"path":       r.URL.Path,
			"proto":      r.Proto,
			"status":     recorder.Code,
			"bytes":      recorder.Body.Len(),
			"referer":    r.Header.Get("referer"),
			"user_agent": r.Header.Get("user-agent"),
			"time_ms":    duration.Milliseconds(),
		}).Info("handled")
	}
	return http.HandlerFunc(fn)
}
