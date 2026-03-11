package grpc

import (
	"context"
	"errors"
	"time"

	userpb "github.com/Hubcher/project-management/contracts/gen/proto/user"
	"github.com/Hubcher/project-management/user-service/internal/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	userpb.UnimplementedUserServiceServer
	service core.UserService
}

func NewServer(service core.UserService) *Server {
	return &Server{service: service}
}

func (s *Server) Ping(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *Server) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.User, error) {
	user, err := s.service.CreateUser(ctx, core.CreateUserInput{
		ID:       req.GetId(),
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Role:     req.GetRole(),
	})
	if err != nil {
		return nil, mapCoreError(err)
	}

	return toProtoUser(user), nil
}

func (s *Server) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.User, error) {
	user, err := s.service.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoUser(user), nil
}

func (s *Server) ListUsers(ctx context.Context, req *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error) {
	users, err := s.service.ListUsers(ctx, req.GetRole())
	if err != nil {
		return nil, mapCoreError(err)
	}

	resp := &userpb.ListUsersResponse{
		Users: make([]*userpb.User, 0, len(users)),
	}

	for i := range users {
		resp.Users = append(resp.Users, toProtoUser(&users[i]))
	}

	return resp, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.User, error) {
	user, err := s.service.UpdateUser(ctx, core.UpdateUserInput{
		ID:       req.GetId(),
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Role:     req.GetRole(),
	})
	if err != nil {
		return nil, mapCoreError(err)
	}

	return toProtoUser(user), nil
}

func (s *Server) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*emptypb.Empty, error) {
	if err := s.service.DeleteUser(ctx, req.GetId()); err != nil {
		return nil, mapCoreError(err)
	}
	return &emptypb.Empty{}, nil
}

func toProtoUser(user *core.User) *userpb.User {
	return &userpb.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
}

func mapCoreError(err error) error {
	switch {
	case errors.Is(err, core.ErrInvalidUser):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, core.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, core.ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
