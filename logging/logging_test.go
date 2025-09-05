package logging

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func FakeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"retval": "done"}`)
}

func TestLoggingMiddleware(t *testing.T) {
	cases := []struct {
		name              string
		loggingMiddleware func(writer io.Writer) func(next http.Handler) http.Handler
	}{
		{
			name: "logrusMiddleware",
			loggingMiddleware: func(writer io.Writer) func(next http.Handler) http.Handler {
				logger := logrus.New()
				logger.SetOutput(writer)
				lr := NewLogrus(WithLogrus(logger))
				return func(next http.Handler) http.Handler {
					return lr.Middleware(next)
				}
			},
		},
		{
			name: "slogMiddleware",
			loggingMiddleware: func(writer io.Writer) func(next http.Handler) http.Handler {
				s := NewSlog(WithSlog(slog.New(slog.NewTextHandler(writer, nil))))
				return func(next http.Handler) http.Handler {
					return s.Middleware(next)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/some-path", nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(FakeHandler)

			// setup logrus
			var capture bytes.Buffer
			writer := bufio.NewWriter(&capture)

			// assert middleware
			tc.loggingMiddleware(writer)(handler).ServeHTTP(rr, req)

			assert.NoError(t, writer.Flush())
			assert.Contains(t, capture.String(), "method=GET path=/some-path proto=HTTP")
			assert.Equal(t, http.StatusOK, rr.Code)
		})
	}
}
