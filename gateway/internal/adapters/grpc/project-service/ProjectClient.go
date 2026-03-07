package project_service

import (
	"log/slog"

	"google.golang.org/grpc"
)

type Client struct {
	log    *slog.Logger
	client projectpb.ProjectServiceClient
	conn   *grpc.ClientConn
}

func NewClient(address string, log *slog.Logger) (*Client, error) {

}
