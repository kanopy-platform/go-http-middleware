package logging

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

func FakeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"retval": "done"}`)
}

func TestLoggingMiddleware(t *testing.T) {
	req, err := http.NewRequest("GET", "/some-path", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(FakeHandler)

	var capture bytes.Buffer

	writer := bufio.NewWriter(&capture)
	log.SetOutput(writer)
	logger := New()

	logger(handler).ServeHTTP(rr, req)
	assert.NoError(t, writer.Flush())
	assert.Contains(t, capture.String(), "method=GET path=/some-path proto=HTTP")
	assert.Equal(t, http.StatusOK, rr.Code)
}
