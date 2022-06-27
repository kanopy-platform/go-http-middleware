package prometheus

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	prom := New(WithDurationBuckets(1, 2, 3))
	assert.Equal(t, []float64{1, 2, 3}, prom.durationBuckets)

	h := prom.Middleware("/", http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	h.ServeHTTP(rr, req)
	h.ServeHTTP(rr, req)

	prom.Handler().ServeHTTP(rr, req)
	body := rr.Body.String()
	assert.Contains(t, body, `http_request_total{code="200",handler="/",method="get"} 2`)
	assert.Contains(t, body, `http_request_duration_seconds_count{code="200",handler="/",method="get"} 2`)
}
