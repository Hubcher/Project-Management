package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	//projectgrpc "github.com/Hubcher/project-management/gateway/internal/adapters/grpc/project"

	usergrpc "github.com/Hubcher/project-management/gateway/internal/adapters/grpc/user"
	"github.com/Hubcher/project-management/gateway/internal/adapters/rest/ping"
	//projectrest "github.com/Hubcher/project-management/gateway/internal/adapters/rest/project"
	userrest "github.com/Hubcher/project-management/gateway/internal/adapters/rest/user"
	"github.com/Hubcher/project-management/gateway/internal/config"
	"github.com/Hubcher/project-management/gateway/internal/core"
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

	log.Info("starting gateway service")
	log.Debug("debug messages are enabled")

	//projectClient, err := projectgrpc.NewClient(cfg.ProjectAddress, log)
	//if err != nil {
	//	log.Error("cannot init ProjectService adapter", "error", err)
	//	return err
	//}
	//defer func() {
	//	if cerr := projectClient.Close(); cerr != nil {
	//		log.Error("cannot close ProjectService adapter", "error", cerr)
	//	}
	//}()

	userClient, err := usergrpc.NewClient(cfg.UserAddress, log)
	if err != nil {
		log.Error("cannot init UserService adapter", "error", err)
		return err
	}
	defer func() {
		if cerr := userClient.Close(); cerr != nil {
			log.Error("cannot close UserService adapter", "error", cerr)
		}
	}()

	mux := http.NewServeMux()

	// userService endpoints
	mux.Handle("POST /api/users", userrest.NewCreateUserHandler(log, userClient))
	mux.Handle("GET /api/users/{id}", userrest.NewGetUserByIdHandler(log, userClient))
	mux.Handle("GET /api/users/{email}", userrest.NewGetUserByEmailHandler(log, userClient))
	mux.Handle("GET /api/users", userrest.NewListUsersHandler(log, userClient))
	mux.Handle("PUT /api/users/{id}", userrest.NewUpdateUserHandler(log, userClient))
	mux.Handle("DELETE /api/users/{id}", userrest.NewDeleteUserHandler(log, userClient))

	// projectService endpoints
	//mux.Handle("POST /api/projects", projectrest.NewCreateProjectHandler(log, projectClient))
	//mux.Handle("GET /api/projects/{contractNumber}", projectrest.NewGetHandler(log, projectClient))
	//mux.Handle("GET /api/projects", projectrest.NewListHandler(log, projectClient))
	//mux.Handle("PUT /api/projects/{contractNumber}", projectrest.NewUpdateHandler(log, projectClient))
	//mux.Handle("DELETE /api/projects/{contractNumber}", projectrest.NewDeleteHandler(log, projectClient))

	// Ping: userService + projectService
	mux.Handle("GET /api/ping", ping.NewPingHandler(log, map[string]core.Pinger{
		//"project": projectClient,
		"user": userClient,
	}))

	server := &http.Server{
		Addr:        cfg.HTTPConfig.Address,
		ReadTimeout: cfg.HTTPConfig.Timeout,
		Handler:     mux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Debug("shutting down server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Error("erroneous shutdown", "error", err)
		}
	}()

	log.Info("running HTTP server", "address", cfg.HTTPConfig.Address)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("server closed unexpectedly", "error", err)
		return err
	}

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
