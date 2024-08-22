package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type contextKey string

const requestIDKey contextKey = "requestId"

var (
	router *mux.Router
)

func getHTTPServer(address, port string) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf("%s:%s", address, port),
		Handler:      router,
		ReadTimeout:  goHTTPReadTimeout,
		WriteTimeout: goHTTPWriteTimeout,
	}
}

func getRouter() {
	router = mux.NewRouter().StrictSlash(true)

	// Apply the logging middleware to all routes
	router.Use(loggingMiddleware)
}

func registerHandlers() {
	for route, handler := range routeMap {
		router.HandleFunc(route, handler)
	}
}

// loggingMiddleware generates a unique request ID, logs the request, and passes the ID to the next handler
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Generate a unique request ID
		requestId := genUniqId()

		// Add it to the response headers
		w.Header().Set("X-Request-ID", requestId)

		// Create a new context with the request ID
		ctx := context.WithValue(r.Context(), requestIDKey, requestId)

		// Create a new http.Request with the new context
		r = r.WithContext(ctx)

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log the request details
		go logInfo(requestId).Str("duration", time.Since(start).String()).Str("method", r.Method).Str("uri", r.RequestURI).Str("remoteAddr", r.RemoteAddr).Msg("")
	})
}

// GetRequestID retrieves the request ID from the context
func getRequestIdFromContext(ctx context.Context) string {
	if requestId, ok := ctx.Value(requestIDKey).(string); ok {
		return requestId
	}
	return "unknown"
}

func setupGracefulShutdown(server *http.Server) {
	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period
		shutdownCtx, _ := context.WithTimeout(serverCtx, goGracefulShutdownTimeout)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				logger.Fatal().Msg("graceful shutdown timed out, forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			logger.Err(err).Msg("failed to shutdown server gracefully")
		}
		serverStopCtx()
	}()

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
