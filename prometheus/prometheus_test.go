package prometheus

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	prom := New()

	h := prom.Middleware("/", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	h.ServeHTTP(rr, req)
	h.ServeHTTP(rr, req)

	prom.Handler().ServeHTTP(rr, req)
	assert.Contains(t, rr.Body.String(), `http_request_total{code="200",handler="/",method="get"} 2`)
	assert.Contains(t, rr.Body.String(), `http_request_duration_seconds_count{code="200",handler="/",method="get"} 2`)
}
