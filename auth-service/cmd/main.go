package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os/signal"
	"syscall"

	authgrpc "github.com/Hubcher/project-management/auth-service/internal/adapters/grpc/auth"
	"github.com/Hubcher/project-management/auth-service/internal/config"
	"github.com/Hubcher/project-management/auth-service/internal/core"
	authpb "github.com/Hubcher/project-management/contracts/gen/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"log/slog"
	"os"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	var configPath string
	flag.StringVar(&configPath, "config", "configs/config.yaml", "path to config.yaml")
	flag.Parse()

	cfg := config.MustLoad(configPath)
	log := mustMakeLogger(cfg.LogLevel, cfg.Env)

	log.Info("starting auth-service")
	log.Debug("debug messages are enabled", slog.Any("cfg", cfg))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	authService := core.NewService(log, db)

	// gRPC server
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	authpb.RegisterAuthServiceServer(s, authgrpc.NewServer(authService))
	reflection.Register(s)

	if err = s.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	go func() {
		<-ctx.Done()
		log.Info("Gracefully shutting down auth-service server")
		s.GracefulStop()
	}()

	return nil
}

func mustMakeLogger(loglevel string, env string) *slog.Logger {
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

	// TODO: сделать для разных env
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}
