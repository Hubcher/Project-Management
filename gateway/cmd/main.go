// Package main provides the public REST API gateway for the project-management system.
//
// @title Project Management Gateway API
// @version 1.0
// @description Public REST API exposed by the gateway. Internal microservices communicate only over gRPC and are intentionally omitted from this specification.
// @host localhost:28080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT access token. Use the `Bearer {token}` format.
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

	authgrpc "github.com/Hubcher/project-management/gateway/internal/adapters/grpc/auth"
	projectgrpc "github.com/Hubcher/project-management/gateway/internal/adapters/grpc/project"
	reportgrpc "github.com/Hubcher/project-management/gateway/internal/adapters/grpc/report"
	usergrpc "github.com/Hubcher/project-management/gateway/internal/adapters/grpc/user"
	authrest "github.com/Hubcher/project-management/gateway/internal/adapters/rest/auth"
	"github.com/Hubcher/project-management/gateway/internal/adapters/rest/middleware"
	"github.com/Hubcher/project-management/gateway/internal/adapters/rest/ping"
	projectrest "github.com/Hubcher/project-management/gateway/internal/adapters/rest/project"
	reportrest "github.com/Hubcher/project-management/gateway/internal/adapters/rest/report"
	swaggerrest "github.com/Hubcher/project-management/gateway/internal/adapters/rest/swagger"
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

	// Auth client
	authClient, err := authgrpc.NewClient(cfg.AuthAddress, log)
	if err != nil {
		log.Error("cannot init AuthService adapter", "error", err)
		return err
	}
	defer func() {
		if cerr := authClient.Close(); cerr != nil {
			log.Error("cannot close AuthService adapter", "error", cerr)
		}
	}()

	// User client
	userClient, err := usergrpc.NewClient(cfg.UserAddress)
	if err != nil {
		log.Error("cannot init UserService adapter", "error", err)
		return err
	}
	defer func() {
		if cerr := userClient.Close(); cerr != nil {
			log.Error("cannot close UserService adapter", "error", cerr)
		}
	}()

	// Project client
	projectClient, err := projectgrpc.NewClient(cfg.ProjectAddress)
	if err != nil {
		log.Error("cannot init ProjectService adapter", "error", err)
		return err
	}
	defer func() {
		if cerr := projectClient.Close(); cerr != nil {
			log.Error("cannot close ProjectService adapter", "error", cerr)
		}
	}()

	// Report client
	reportClient, err := reportgrpc.NewClient(cfg.ReportAddress)
	if err != nil {
		log.Error("cannot init ReportService adapter", "error", err)
		return err
	}
	defer func() {
		if cerr := reportClient.Close(); cerr != nil {
			log.Error("cannot close ReportService adapter", "error", cerr)
		}
	}()

	identity := core.NewIdentityService(authClient, userClient)
	authMW := middleware.NewAuthMiddleware(authClient)
	mux := http.NewServeMux()

	// swagger endpoints
	mux.Handle("GET /swagger/", swaggerrest.NewUIHandler())
	mux.Handle("GET /openapi/", swaggerrest.NewSpecHandler())

	// auth endpoints
	mux.Handle("POST /api/auth/register", authrest.NewRegisterHandler(identity))
	mux.Handle("POST /api/auth/login", authrest.NewLoginHandler(identity))
	mux.Handle("GET /api/auth/me", middleware.Chain(authrest.NewMeHandler(), authMW.Auth))

	// user endpoints
	mux.Handle("POST /api/users", middleware.Chain(userrest.NewCreateUserHandler(identity), authMW.Auth, authMW.RequireRoles(core.RoleAdmin)))
	mux.Handle("GET /api/users/{id}", middleware.Chain(userrest.NewGetUserByIDHandler(userClient), authMW.Auth))
	mux.Handle("GET /api/users", middleware.Chain(userrest.NewListUsersHandler(userClient), authMW.Auth, authMW.RequireRoles(core.RoleAdmin)))
	mux.Handle("PUT /api/users/{id}", middleware.Chain(userrest.NewUpdateUserHandler(userClient), authMW.Auth))
	mux.Handle("DELETE /api/users/{id}", middleware.Chain(userrest.NewDeleteUserHandler(identity), authMW.Auth, authMW.RequireRoles(core.RoleAdmin)))

	// project endpoints
	mux.Handle("POST /api/projects", middleware.Chain(projectrest.NewCreateProjectHandler(projectClient, userClient), authMW.Auth))
	mux.Handle("GET /api/projects/{id}", middleware.Chain(projectrest.NewGetProjectHandler(projectClient), authMW.Auth))
	mux.Handle("GET /api/projects", middleware.Chain(projectrest.NewListProjectsHandler(projectClient), authMW.Auth))
	mux.Handle("PUT /api/projects/{id}", middleware.Chain(projectrest.NewUpdateProjectHandler(projectClient, userClient), authMW.Auth))
	mux.Handle("DELETE /api/projects/{id}", middleware.Chain(projectrest.NewDeleteProjectHandler(projectClient), authMW.Auth))

	// project stages endpoints
	mux.Handle("POST /api/projects/{projectId}/stages", middleware.Chain(projectrest.NewCreateStageHandler(projectClient), authMW.Auth))
	mux.Handle("GET /api/projects/{projectId}/stages", middleware.Chain(projectrest.NewListStagesHandler(projectClient), authMW.Auth))
	mux.Handle("GET /api/project-stages/{id}", middleware.Chain(projectrest.NewGetStageHandler(projectClient), authMW.Auth))
	mux.Handle("PUT /api/project-stages/{id}", middleware.Chain(projectrest.NewUpdateStageHandler(projectClient), authMW.Auth))
	mux.Handle("DELETE /api/project-stages/{id}", middleware.Chain(projectrest.NewDeleteStageHandler(projectClient), authMW.Auth))

	// project members endpoints
	mux.Handle("POST /api/projects/{projectId}/members", middleware.Chain(projectrest.NewCreateMemberHandler(projectClient, userClient), authMW.Auth))
	mux.Handle("GET /api/projects/{projectId}/members", middleware.Chain(projectrest.NewListMembersHandler(projectClient), authMW.Auth))
	mux.Handle("GET /api/project-members/{id}", middleware.Chain(projectrest.NewGetMemberHandler(projectClient), authMW.Auth))
	mux.Handle("PUT /api/project-members/{id}", middleware.Chain(projectrest.NewUpdateMemberHandler(projectClient), authMW.Auth))
	mux.Handle("DELETE /api/project-members/{id}", middleware.Chain(projectrest.NewDeleteMemberHandler(projectClient), authMW.Auth))

	// project events endpoints
	mux.Handle("POST /api/projects/{projectId}/events", middleware.Chain(projectrest.NewCreateEventHandler(projectClient), authMW.Auth))
	mux.Handle("GET /api/projects/{projectId}/events", middleware.Chain(projectrest.NewListEventsHandler(projectClient), authMW.Auth))
	mux.Handle("GET /api/project-events/{id}", middleware.Chain(projectrest.NewGetEventHandler(projectClient), authMW.Auth))
	mux.Handle("PUT /api/project-events/{id}", middleware.Chain(projectrest.NewUpdateEventHandler(projectClient), authMW.Auth))
	mux.Handle("DELETE /api/project-events/{id}", middleware.Chain(projectrest.NewDeleteEventHandler(projectClient), authMW.Auth))

	// report endpoints
	mux.Handle("POST /api/reports", middleware.Chain(reportrest.NewCreateReportHandler(reportClient), authMW.Auth))
	mux.Handle("GET /api/reports/{id}", middleware.Chain(reportrest.NewGetReportHandler(reportClient, projectClient), authMW.Auth))
	mux.Handle("GET /api/reports", middleware.Chain(reportrest.NewListReportsHandler(reportClient, projectClient), authMW.Auth))
	mux.Handle("PUT /api/reports/{id}", middleware.Chain(reportrest.NewUpdateReportHandler(reportClient, projectClient), authMW.Auth))
	mux.Handle("DELETE /api/reports/{id}", middleware.Chain(reportrest.NewDeleteReportHandler(reportClient), authMW.Auth))

	// report-entries endpoints
	mux.Handle("POST /api/reports/{reportId}/entries", middleware.Chain(reportrest.NewCreateEntryHandler(reportClient, projectClient), authMW.Auth))
	mux.Handle("GET /api/reports/{reportId}/entries", middleware.Chain(reportrest.NewListEntriesHandler(reportClient, projectClient), authMW.Auth))
	mux.Handle("GET /api/report-entries/{id}", middleware.Chain(reportrest.NewGetEntryHandler(reportClient, projectClient), authMW.Auth))
	mux.Handle("PUT /api/report-entries/{id}", middleware.Chain(reportrest.NewUpdateEntryHandler(reportClient, projectClient), authMW.Auth))
	mux.Handle("DELETE /api/report-entries/{id}", middleware.Chain(reportrest.NewDeleteEntryHandler(reportClient, projectClient), authMW.Auth))

	// report comment endpoints
	mux.Handle("POST /api/reports/{reportId}/comments", middleware.Chain(reportrest.NewCreateCommentHandler(reportClient, projectClient), authMW.Auth))
	mux.Handle("GET /api/reports/{reportId}/comments", middleware.Chain(reportrest.NewListCommentsHandler(reportClient, projectClient), authMW.Auth))
	mux.Handle("GET /api/report-comments/{id}", middleware.Chain(reportrest.NewGetCommentHandler(reportClient, projectClient), authMW.Auth))
	mux.Handle("PUT /api/report-comments/{id}", middleware.Chain(reportrest.NewUpdateCommentHandler(reportClient, projectClient), authMW.Auth))
	mux.Handle("DELETE /api/report-comments/{id}", middleware.Chain(reportrest.NewDeleteCommentHandler(reportClient, projectClient), authMW.Auth))

	// system ping endpoint
	mux.Handle("GET /api/ping", ping.NewPingHandler(log,
		map[string]core.Pinger{
			"auth":    authClient,
			"user":    userClient,
			"project": projectClient,
			"report":  reportClient,
		}))

	server := &http.Server{Addr: cfg.HTTPConfig.Address, ReadTimeout: cfg.HTTPConfig.Timeout, Handler: mux}
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
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
}
