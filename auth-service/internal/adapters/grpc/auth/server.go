package auth

import (
    "context"
    "errors"
    "strings"

    "github.com/Hubcher/project-management/auth-service/internal/core"
    authpb "github.com/Hubcher/project-management/contracts/gen/go/auth"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/metadata"
    "google.golang.org/grpc/status"
    "google.golang.org/protobuf/types/known/emptypb"
)

const accountRoleMetadataKey = "x-account-role"

type Server struct {
    authpb.UnimplementedAuthServiceServer
    service *core.Service
}

func NewServer(svc *core.Service) *Server {
    return &Server{service: svc}
}

func (s *Server) Ping(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
    return &emptypb.Empty{}, nil
}

func (s *Server) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
    role := extractRequestedRole(ctx)

    userID, err := s.service.RegisterWithRole(ctx, req.GetEmail(), req.GetPassword(), role)
    if err != nil {
        return &authpb.RegisterResponse{
            Status: 500,
            Error:  err.Error(),
        }, mapCoreError(err)
    }
    return &authpb.RegisterResponse{
        UserId: userID,
        Status: 201,
    }, nil
}

func (s *Server) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
    token, err := s.service.Login(ctx, req.GetEmail(), req.GetPassword())
    if err != nil {
        return &authpb.LoginResponse{
            Status: 401,
            Error:  err.Error(),
        }, mapCoreError(err)
    }
    return &authpb.LoginResponse{
        Token:  token,
        Status: 200,
    }, nil
}

func (s *Server) Validate(ctx context.Context, req *authpb.ValidateRequest) (*authpb.ValidateResponse, error) {
    claims, err := s.service.Validate(ctx, req.GetToken())
    if err != nil {
        return &authpb.ValidateResponse{
            Status: 401,
            Error:  err.Error(),
        }, mapCoreError(err)
    }
    return &authpb.ValidateResponse{
        Status: 200,
        UserId: claims.UserID,
        Role:   string(claims.Role),
        Email:  claims.Email,
    }, nil
}

func (s *Server) DeleteCredentials(ctx context.Context, req *authpb.DeleteCredentialsRequest) (*authpb.DeleteCredentialsResponse, error) {
    if err := s.service.DeleteCredentials(ctx, req.GetUserID()); err != nil {
        return &authpb.DeleteCredentialsResponse{
            Status: 400,
            Error:  err.Error(),
        }, mapCoreError(err)
    }
    return &authpb.DeleteCredentialsResponse{
        Status: 200,
    }, nil
}

func extractRequestedRole(ctx context.Context) core.Role {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return core.RoleUser
    }

    values := md.Get(accountRoleMetadataKey)
    if len(values) == 0 {
        return core.RoleUser
    }

    switch strings.ToLower(strings.TrimSpace(values[0])) {
    case string(core.RoleAdmin):
        return core.RoleAdmin
    default:
        return core.RoleUser
    }
}

func mapCoreError(err error) error {
    switch {
    case errors.Is(err, core.ErrInvalidArgument):
        return status.Error(codes.InvalidArgument, err.Error())
    case errors.Is(err, core.ErrAccountExists):
        return status.Error(codes.AlreadyExists, err.Error())
    case errors.Is(err, core.ErrInvalidCredentials):
        return status.Error(codes.Unauthenticated, err.Error())
    case errors.Is(err, core.ErrInvalidToken):
        return status.Error(codes.Unauthenticated, err.Error())
    case errors.Is(err, core.ErrForbidden):
        return status.Error(codes.PermissionDenied, err.Error())
    case errors.Is(err, core.ErrInactiveAccount):
        return status.Error(codes.PermissionDenied, err.Error())
    case errors.Is(err, core.ErrAccountNotFound):
        return status.Error(codes.NotFound, err.Error())
    default:
        return status.Error(codes.Internal, "internal error")
    }
}
