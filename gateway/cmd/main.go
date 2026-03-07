package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"project-managment/gateway/config"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)
	log := mustMakeLogger(cfg.LogLevel)

	log.Info("starting gateway service")
	log.Debug("debug messages are enabled")

	// Инициализируем мультиплексер
	mux := http.NewServeMux()

	// project-service CRUD

	return nil
}

func mustMakeLogger(loglevel string) *slog.Logger {
	var level slog.Level

	switch loglevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		panic("invalid log level " + loglevel)
	}

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}
