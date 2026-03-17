package user

import (
	"log/slog"

	userpb "github.com/Hubcher/project-management/contracts/gen/proto/user"
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client userpb.userServiceClient
	log    *slog.Logger
}
