package prometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Prometheus struct {
	duration *prometheus.HistogramVec
	counter  *prometheus.CounterVec
}

type OptionFunc func(*Prometheus)

func New(opts ...OptionFunc) *Prometheus {
	labelNames := []string{"code", "handler", "method"}

	m := &Prometheus{
		counter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_request_total",
				Help: "A counter for http requests.",
			},
			labelNames,
		),
		duration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "A histogram of latencies for http requests.",
				Buckets: []float64{.25, .5, 1, 2.5, 5, 10},
			},
			labelNames,
		),
	}

	for _, opt := range opts {
		opt(m)
	}

	prometheus.MustRegister(m.counter)
	prometheus.MustRegister(m.duration)

	return m
}

func (m *Prometheus) Middleware(path string, next http.Handler) http.Handler {
	labels := prometheus.Labels{"handler": path}

	return promhttp.InstrumentHandlerDuration(m.duration.MustCurryWith(labels),
		promhttp.InstrumentHandlerCounter(m.counter.MustCurryWith(labels), next))
}

func (m *Prometheus) Handler() http.Handler {
	return promhttp.Handler()
}
