package main

import (
	"context"
	"flag"
	"fmt"

	projectpb "github.com/Hubcher/project-management/contracts/gen/proto/project"
	projectgrpc "github.com/Hubcher/project-management/project-service/internal/adapters/grpc"
	"github.com/Hubcher/project-management/project-service/internal/config"
	"github.com/Hubcher/project-management/project-service/internal/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
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

	log.Info("starting project-service server")
	log.Debug("debug message are enabled")

	projectService := core.NewService(log)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// gRPC server
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	projectpb.RegisterProjectServiceServer(s, projectgrpc.NewServer(projectService))
	reflection.Register(s)

	go func() {
		<-ctx.Done()
		log.Info("Gracefully shutting down project-service server")
	}()

	if err = s.Serve(listener); err != nil {
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

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}
