package auth

import (
	"context"

	"github.com/Hubcher/project-management/auth-service/internal/core"
	authpb "github.com/Hubcher/project-management/contracts/gen/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/protobuf/types/known/emptypb"
)

const emptyValue = 0

type Server struct {
	authpb.UnimplementedAuthServiceServer
	service core.Service
	auth    Auth // TODO: убрать в сервис?
}

func NewServer(service core.Service) *Server {
	return &Server{service: service}
}

func (s *Server) Ping(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

type Auth interface {
	Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (UserID string, err error)
	IsAdmin(ctx context.Context, email string) (bool, error)
}

func (s *Server) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	// TODO: через какой-нибудь package сделать валидацию и вынести в отдельную функцию
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		
		return nil, status.Error(codes.InvalidArgument, "internal error")
	}

	return &authpb.LoginResponse{Token: token}, nil
}

func (s *Server) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		// TODO: .... ошибки по типу уже существующего пользователя
		return nil, status.Error(codes.InvalidArgument, "internal error")
	}

	return &authpb.RegisterResponse{UserId: userID}, nil

}

func (s *Server) IsAdmin(ctx context.Context, req *authpb.IsAdminRequest) (*authpb.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.UserId)
	if err != nil {
		// TODO: ...
		return nil, status.Error(codes.InvalidArgument, "internal error")
	}
	return &authpb.IsAdminResponse{IsAdmin: isAdmin}, nil
}

func validateLogin(req *authpb.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}
	return nil
}

func validateRegister(req *authpb.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	return nil
}

func validateIsAdmin(req *authpb.IsAdminRequest) error {
	if req.GetUserId() == "" {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}
	return nil
}
