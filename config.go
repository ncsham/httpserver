package main

const (
	PROJECT = "httpServer"
)

var (
	address = getEnv("ADDRESS", "0.0.0.0")
	port    = getEnv("PORT", "20000")

	httpReadTimeout         = getEnv("HTTP_READ_TIMEOUT", "10s")
	httpWriteTimeout        = getEnv("HTTP_WRITE_TIMEOUT", "10s")
	gracefulShutdownTimeout = getEnv("GRACEFUL_SHUTDOWN_TIMEOUT", "10s")
)
