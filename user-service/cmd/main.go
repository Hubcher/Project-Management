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

	userpb "github.com/Hubcher/project-management/contracts/gen/proto/user"
	"github.com/Hubcher/project-management/user-service/internal/adapters/db/postgres"
	usergrpc "github.com/Hubcher/project-management/user-service/internal/adapters/grpc"
	"github.com/Hubcher/project-management/user-service/internal/config"
	"github.com/Hubcher/project-management/user-service/internal/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	var configPath string
	flag.StringVar(&configPath, "config", "configs/config.yaml", "Path to config file")
	flag.Parse()

	cfg := config.MustLoad(configPath)
	log := mustMakeLogger(cfg.LogLevel)

	log.Info("starting user-service server")
	log.Debug("debug message are enabled")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// database postgres adapter
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

	userService := core.NewService(log, storage)

	// gRPC server
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	userpb.RegisterUserServiceServer(s, usergrpc.NewServer(userService))
	reflection.Register(s)

	go func() {
		<-ctx.Done()
		log.Info("Gracefully shutting down user-service server")
		s.GracefulStop()
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
