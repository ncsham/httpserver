package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func initPrometheus() {
	// Create a new ServeMux for Prometheus metrics
	metricsMux := http.NewServeMux()

	// Register the Prometheus handler at the desired path using the dedicated mux
	metricsMux.Handle(promMetricsPath, promhttp.Handler())

	// Create a new server for metrics
	metricsServer := &http.Server{
		Addr:    ":" + promMetricsPort,
		Handler: metricsMux,
	}

	// Start the metrics server in a goroutine
	go func() {
		err := metricsServer.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				// The server was closed, which is expected when the application is shutting down
				return
			}
			// Handle the error, e.g., log it or notify the user
			logger.Panic().Stack().Err(err).Msg("Error starting Prometheus server")
		}
	}()

	logger.Info().Str("port", promMetricsPort).Str("path", "/metrics").Msg("Prometheus server started")
}

var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Count of all HTTP requests",
	}, []string{"code", "method", "path"})

	httpRequestsSummary = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "http_request_duration_seconds",
		Help:       "Histogram of request latencies",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"path"})
)
