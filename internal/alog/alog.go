package alog

import (
	"io"
	"log/slog"
	"os"
)

func New(logFilename string, devMode bool) (*slog.Logger, func()) {
	w := logWriter(logFilename, devMode)
	loggerOpts := &slog.HandlerOptions{}
	if devMode {
		loggerOpts.Level = slog.LevelDebug
	}
	logger := slog.New(logHandler(w, loggerOpts, devMode))

	slog.SetDefault(logger)

	return logger, func() {
		if f, ok := w.(*os.File); ok {
			f.Close()
		}
	}
}

func Error(err error) slog.Attr {
	return slog.String("error", err.Error())
}

func logHandler(w io.Writer, opts *slog.HandlerOptions, devMode bool) slog.Handler {
	if devMode {
		return slog.NewTextHandler(w, opts)
	}
	return slog.NewJSONHandler(w, opts)
}

func logWriter(logFilename string, devMode bool) io.Writer {
	if logFilename == "" || devMode {
		return os.Stdout
	}

	f, err := os.OpenFile(logFilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	return f
}
