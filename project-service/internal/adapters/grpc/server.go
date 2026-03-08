package grpc

import (
	"context"

	projectpb "github.com/Hubcher/project-management/contracts/gen/proto/project"
	"github.com/Hubcher/project-management/project-service/internal/core"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	projectpb.UnimplementedProjectServiceServer
	service core.ProjectRepo
}

func NewServer(service core.ProjectRepo) *Server {
	return &Server{service: service}
}

func (s *Server) Ping(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
