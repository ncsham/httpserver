package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"
)

func getCmdLineArgs() {
	flag.BoolVar(&debugProgram, "debug", false, "Debug")
	flag.Parse()
}

func getStartUpConf() {
	var err error

	if !debugProgram {
		if debugProgram, err = strconv.ParseBool(getEnv("DEBUG", "false")); err != nil {
			debugProgram = false
		}
	}

	hostname = getHostname()

	if goGracefulShutdownTimeout, err = time.ParseDuration(gracefulShutdownTimeout); err != nil {
		fmt.Println(fmt.Sprintf("Error parsing graceful shutdown timeout: %s", err.Error()))
		goGracefulShutdownTimeout = time.Duration(10) * time.Second
	}

	if goHTTPReadTimeout, err = time.ParseDuration(httpReadTimeout); err != nil {
		fmt.Println(fmt.Sprintf("Error parsing http read timeout: %s", err.Error()))
		goHTTPReadTimeout = time.Duration(10) * time.Second
	}

	if goHTTPWriteTimeout, err = time.ParseDuration(httpWriteTimeout); err != nil {
		fmt.Println(fmt.Sprintf("Error parsing http write timeout: %s", err.Error()))
		goHTTPWriteTimeout = time.Duration(10) * time.Second
	}

	startedAt = time.Now().Format(time.RFC3339)

	fmt.Println(fmt.Sprintf("Starting Program with debug: %t, hostname: %s, startedAt: %s", debugProgram, hostname, startedAt))
}
