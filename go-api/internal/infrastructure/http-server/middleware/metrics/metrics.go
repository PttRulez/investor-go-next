package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	// http method of request
	labelMethod = "method"

	// path of request
	labelPath = "path"

	// http code returned from our server
	labelCode = "code"
)

type Metrics struct {
	reqTotal         prometheus.Counter
	reqLatency       *prometheus.HistogramVec
	errInternalTotal prometheus.Counter
}

func Register() *Metrics {
	reqTotal := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "http",
		Name:      "req_total",
		Help:      "Total number of requests received.",
	})
	reqLatency := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "http",
			Name:      "latency",
			Help:      "Duration of request in seconds.",
			Buckets:   []float64{0.1, 0.5, 1},
		},
		[]string{labelMethod, labelPath, labelCode},
	)
	errInternalTotal := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "http",
		Name:      "500_errors",
		Help:      "Total number of requests received.",
	})

	return &Metrics{
		reqTotal:         reqTotal,
		reqLatency:       reqLatency,
		errInternalTotal: errInternalTotal,
	}
}

func (m *Metrics) Default() http.Handler {
	return promhttp.Handler()
}

// Increments Requests Counter
func (m *Metrics) IncTotalRequests() {
	m.reqTotal.Inc()
}

// Observes Request Duration (Latency)
func (m *Metrics) WriteLatency(method, path string, duration float64) {
	m.reqLatency.WithLabelValues(method, path).Observe(duration)
}

// Increments 500 Errors Counter
func (m *Metrics) IncTotalInternalErrors() {
	m.errInternalTotal.Inc()
}
