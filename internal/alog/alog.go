package alog

import (
	"io"
	"log/slog"
	"os"
)

func New(logFilename string, devMode bool) (*slog.Logger, func()) {
	var w io.Writer = os.Stdout

	fw := fileWriter(logFilename, devMode)
	if fw != nil {
		w = io.MultiWriter(os.Stdout, fw)
	}

	loggerOpts := &slog.HandlerOptions{}
	if devMode {
		loggerOpts.Level = slog.LevelDebug
	}
	logger := slog.New(logHandler(w, loggerOpts, devMode))

	slog.SetDefault(logger)

	return logger, func() {
		if fw != nil {
			fw.Close()
		}
	}
}

func logHandler(w io.Writer, opts *slog.HandlerOptions, devMode bool) slog.Handler {
	if devMode {
		return slog.NewTextHandler(w, opts)
	}
	return slog.NewJSONHandler(w, opts)
}

func fileWriter(logFilename string, devMode bool) *os.File {
	if logFilename == "" || devMode {
		return nil
	}

	f, err := os.OpenFile(logFilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	return f
}

func Error(err error) slog.Attr {
	return slog.String("error", err.Error())
}
