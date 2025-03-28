package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

var routeMap = map[string]func(w http.ResponseWriter, req *http.Request){
	"/status": statusHandler,
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, fmt.Sprintf("Served with RequestID: %s", getRequestIdFromContext(r.Context())))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, "Dont Try Again")

	duration := time.Since(start)

	// Record metrics with "NA" for path
	httpRequestsTotal.WithLabelValues(
		strconv.Itoa(http.StatusNotFound),
		r.Method,
		"NA",
	).Inc()

	httpRequestsSummary.WithLabelValues("NA").Observe(duration.Seconds())

	// Log the not found request
	go logger.Info().
		Str("duration", duration.String()).
		Str("method", r.Method).
		Str("uri", r.RequestURI).
		Str("remoteAddr", r.RemoteAddr).
		Int("status", http.StatusNotFound).
		Msg("Not Found")
}
