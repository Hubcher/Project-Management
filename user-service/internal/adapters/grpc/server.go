package grpc

import (
	userpb "github.com/Hubcher/project-management/contracts/gen/proto/user"
	"github.com/Hubcher/project-management/user-service/internal/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// TODO:
//func Server struct {
//	userpb.UnimplementedUserServiceServer
//	service core.UserService
//}
