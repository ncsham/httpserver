package main

import (
	"net/http"
	"time"
)

var (
	debugProgram bool
	hostname     string

	goHTTPReadTimeout         time.Duration
	goHTTPWriteTimeout        time.Duration
	goGracefulShutdownTimeout time.Duration

	startedAt              string
	packageVersion         string
	packageCommit          string
	packageCommitTimestamp string
	packageBuildTime       string
)

func main() {
	getCmdLineArgs()
	getStartUpConf()
	setLogging()

	getRouter()
	registerHandlers()
	server := getHTTPServer(address, port)

	// Setup graceful shutdown
	go setupGracefulShutdown(server)

	logInfo("main").Str("hostname", hostname).Str("packageVersion", packageVersion).Str("packageCommit", packageCommit).Str("packageBuildTime", packageBuildTime).Str("packageCommitTimestamp", packageCommitTimestamp).Str("startedAt", startedAt).Msg("Starting Server")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logError("main").Err(err).Msg("Error starting server")
	}

	logger.Info().Str("hostname", hostname).Str("StoppedAt", time.Now().Format(time.RFC3339)).Msg("Server Stopped")
}
