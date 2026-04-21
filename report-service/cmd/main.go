package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	reportpb "github.com/Hubcher/project-management/contracts/gen/go/report"
	"github.com/Hubcher/project-management/report-service/internal/adapters/db/postgres"
	reportgrpc "github.com/Hubcher/project-management/report-service/internal/adapters/grpc"
	"github.com/Hubcher/project-management/report-service/internal/config"
	"github.com/Hubcher/project-management/report-service/internal/core"
	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	var configPath string
	flag.StringVar(&configPath, "config", "configs/config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)
	log := mustMakeLogger(cfg.LogLevel)

	log.Info("starting report-service server")
	log.Debug("debug messages are enabled", slog.Any("cfg", cfg))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	storage, err := postgres.New(log, cfg.DBAddress)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer func() {
		if cerr := storage.Close(); cerr != nil {
			log.Error("failed to close database connection", "error", cerr)
		}
	}()

	if err = storage.Migrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	reportService := core.NewService(storage)

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	reportpb.RegisterReportServiceServer(server, reportgrpc.NewServer(reportService))

	go func() {
		<-ctx.Done()
		log.Info("gracefully shutting down report-service server")
		server.GracefulStop()
	}()

	if err = server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

func mustMakeLogger(logLevel string) *slog.Logger {
	var level slog.Level

	switch logLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "ERROR":
		level = slog.LevelError
	case "WARN":
		level = slog.LevelWarn
	default:
		panic("unknown log level: " + logLevel)
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}