package alog

import (
	"io"
	"log/slog"
	"os"
)

var logLevel = new(slog.LevelVar)

func New(devMode bool) *slog.Logger {
	loggerOpts := &slog.HandlerOptions{Level: logLevel}
	logger := slog.New(logHandler(os.Stdout, loggerOpts, devMode))
	if devMode {
		logLevel.Set(slog.LevelDebug)
	}
	slog.SetDefault(logger)

	return logger
}

func SetLevel(level slog.Level) {
	logLevel.Set(level)
}

func logHandler(w io.Writer, opts *slog.HandlerOptions, devMode bool) slog.Handler {
	if devMode {
		return slog.NewTextHandler(w, opts)
	}
	return slog.NewJSONHandler(w, opts)
}

func Error(err error) slog.Attr {
	return slog.String("error", err.Error())
}
