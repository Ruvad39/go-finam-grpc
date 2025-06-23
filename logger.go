package finam

import (
	"log/slog"
	"os"
)

var logLevel = &slog.LevelVar{} // INFO

var log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	Level: logLevel,
})).With(slog.String("package", "go-finam-grpc"))

func SetLogger(logger *slog.Logger) {
	log = logger
}

// SetLogLevel установим уровень логирования
func SetLogLevel(level slog.Level) {
	logLevel.Set(level)
}

// SetLogDebug установим уровень логгирования Debug
func SetLogDebug(debug bool) {
	if debug {
		logLevel.Set(slog.LevelDebug)
	} else {
		logLevel.Set(slog.LevelInfo)
	}
}
