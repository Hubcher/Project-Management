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

	paymentpb "github.com/Hubcher/project-management/contracts/gen/go/paymentcalendar"
	"github.com/Hubcher/project-management/payment-calendar-service/internal/adapters/db/postgres"
	paymentgrpc "github.com/Hubcher/project-management/payment-calendar-service/internal/adapters/grpc"
	projectgrpc "github.com/Hubcher/project-management/payment-calendar-service/internal/adapters/grpc/project"
	"github.com/Hubcher/project-management/payment-calendar-service/internal/config"
	"github.com/Hubcher/project-management/payment-calendar-service/internal/core"
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

	log.Info("starting payment-calendar-service server")
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

	projectClient, err := projectgrpc.NewClient(cfg.ProjectAddress)
	if err != nil {
		return fmt.Errorf("failed to connect to project service: %v", err)
	}
	defer func() {
		if cerr := projectClient.Close(); cerr != nil {
			log.Error("failed to close project service connection", "error", cerr)
		}
	}()

	paymentService := core.NewService(storage, projectClient)

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	paymentpb.RegisterPaymentCalendarServiceServer(server, paymentgrpc.NewServer(paymentService))

	go func() {
		<-ctx.Done()
		log.Info("gracefully shutting down payment-calendar-service server")
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
