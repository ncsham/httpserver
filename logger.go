package main

import (
	"io"
	"os"
	"path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger    *Logger
	logs_dir  string
	logs_file string
)

// Configuration for logging
type LoggingConfig struct {
	// Enable Debug Logging
	DebugLoggingEnabled bool
	// Enable console logging
	ConsoleLoggingEnabled bool
	// EncodeLogsAsJson makes the log framework log JSON
	EncodeLogsAsJson bool
	// FileLoggingEnabled makes the framework log to a file
	// the fields below can be skipped if this value is false!
	FileLoggingEnabled bool
	// Directory to log to to when filelogging is enabled
	Directory string
	// Filename is the name of the logfile which will be placed inside the directory
	Filename string
	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int
	// MaxBackups the max number of rolled files to keep
	MaxBackups int
	// MaxAge the max age in days to keep a logfile
	MaxAge int
}

type Logger struct {
	*zerolog.Logger
}

// Configure sets up the logging framework
//
// In production, the container logs will be collected and file logging should be disabled. However,
// during development it's nicer to see logs as text and optionally write to a file when debugging
// problems in the containerized pipeline
//
// The output log file will be located at /var/log/service-xyz/service-xyz.log and
// will be rolled according to configuration set.
func Configure(config LoggingConfig) *Logger {
	var writers []io.Writer

	if config.ConsoleLoggingEnabled {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}
	if config.FileLoggingEnabled {
		writers = append(writers, newRollingFile(config))
	}

	mw := io.MultiWriter(writers...)

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debugProgram {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Configure Error Stack Marshaler for getting formatted stacktrace
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	logger := zerolog.New(mw).With().Timestamp().Logger()

	logger.Info().
		Bool("jsonLogOutput", config.EncodeLogsAsJson).
		Msg("logging configured")

	return &Logger{
		Logger: &logger,
	}
}

func newRollingFile(config LoggingConfig) io.Writer {
	if err := os.MkdirAll(config.Directory, 0744); err != nil {
		log.Error().Err(err).Str("path", config.Directory).Msg("can't create log directory")
		return nil
	}

	return &lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxBackups: config.MaxBackups, // files
		MaxSize:    config.MaxSize,    // megabytes
		MaxAge:     config.MaxAge,     // days
		Compress:   true,
	}
}

func setLogging() {
	var logConfig LoggingConfig
	logConfig.DebugLoggingEnabled = debugProgram
	logConfig.ConsoleLoggingEnabled = true
	logConfig.EncodeLogsAsJson = true
	logConfig.FileLoggingEnabled = false
	// Directory to log to to when filelogging is enabled
	logConfig.Directory = logs_dir
	// Filename is the name of the logfile which will be placed inside the directory
	logConfig.Filename = logs_file
	// MaxSize the max size in MB of the logfile before it's rolled
	logConfig.MaxSize = 512
	// MaxBackups the max number of rolled files to keep
	logConfig.MaxBackups = 10
	// MaxAge the max age in days to keep a logfile
	logConfig.MaxAge = 10

	logger = Configure(logConfig)
}

func logDebug(rqid string) *zerolog.Event {
	return logger.Debug().Str("rqid", rqid)
}

func logInfo(rqid string) *zerolog.Event {
	return logger.Info().Str("rqid", rqid)
}

func logWarning(rqid string) *zerolog.Event {
	return logger.Warn().Str("rqid", rqid)
}

func logError(rqid string) *zerolog.Event {
	return logger.Error().Stack().Str("rqid", rqid)
}
