package user

import (
	"context"
	"errors"
	"log/slog"
	"time"

	userpb "github.com/Hubcher/project-management/contracts/gen/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	log    *slog.Logger
	client userpb.UserServiceClient
	conn   *grpc.ClientConn
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  1 * time.Second,
				Multiplier: 1.6,
				MaxDelay:   5 * time.Second,
			},
			MinConnectTimeout: 10 * time.Second,
		}),
	)
	if err != nil {
		return nil, err
	}
	conn.Connect()

	return &Client{
		log:    log,
		client: userpb.NewUserServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.client.Ping(ctx, &emptypb.Empty{})
	return err
}

func (c *Client) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.User, error) {
	if req == nil {
		return nil, errors.New("create user request is empty")
	}

	resp, err := c.client.CreateUser(ctx, req)
	if err != nil {
		c.log.Error("CreateUser failed", "error", err)
		return nil, err
	}

	return resp, nil
}

func (c *Client) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.User, error) {
	if req == nil {
		return nil, errors.New("get user request is empty")
	}
	resp, err := c.client.GetUser(ctx, req)
	if err != nil {
		c.log.Error("GetUser failed", "error", err)
		return nil, err
	}

	return resp, nil
}

func (c *Client) ListUsers(ctx context.Context, req *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error) {
	if req == nil {
		return nil, errors.New("list users request is empty")
	}
	resp, err := c.client.ListUsers(ctx, req)
	if err != nil {
		c.log.Error("ListUsers failed", "error", err)
		return nil, err
	}

	return resp, nil
}

func (c *Client) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.User, error) {
	if req == nil {
		return nil, errors.New("update user request is empty")
	}
	resp, err := c.client.UpdateUser(ctx, req)
	if err != nil {
		c.log.Error("UpdateUser failed", "error", err)
		return nil, err
	}
	return resp, nil
}

func (c *Client) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, errors.New("delete user request is empty")
	}
	resp, err := c.client.DeleteUser(ctx, req)
	if err != nil {
		c.log.Error("DeleteUser failed", "error", err)
		return nil, err
	}
	return resp, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
