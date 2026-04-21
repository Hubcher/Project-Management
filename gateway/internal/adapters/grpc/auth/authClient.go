package auth

import (
	"context"
	"log/slog"
	"time"

	authpb "github.com/Hubcher/project-management/contracts/gen/go/auth"
	"github.com/Hubcher/project-management/gateway/internal/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

const accountRoleMetadataKey = "x-account-role"

type Client struct {
	log    *slog.Logger
	conn   *grpc.ClientConn
	client authpb.AuthServiceClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  1 * time.Second,
				Multiplier: 1.6,
				MaxDelay:   10 * time.Second,
			},
			MinConnectTimeout: 10 * time.Second,
		}),
	)
	if err != nil {
		return nil, err
	}
	conn.Connect()

	return &Client{
		client: authpb.NewAuthServiceClient(conn),
		log:    log,
		conn:   conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.client.Ping(ctx, &emptypb.Empty{})
	return err
}

func (c *Client) Register(ctx context.Context, email, password string, role core.Role) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if role != "" {
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(accountRoleMetadataKey, string(role)))
	}

	resp, err := c.client.Register(ctx, &authpb.RegisterRequest{
		Email:    email,
		Password: password,
	})

	if err != nil {
		return "", err
	}
	return resp.GetUserId(), nil
}

func (c *Client) Login(ctx context.Context, email, password string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.Login(ctx, &authpb.LoginRequest{
		Email:    email,
		Password: password,
	})

	if err != nil {
		return "", err
	}
	return resp.GetToken(), nil
}

func (c *Client) Validate(ctx context.Context, token string) (core.AuthUser, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.Validate(ctx, &authpb.ValidateRequest{Token: token})
	if err != nil {
		return core.AuthUser{}, err
	}

	return core.AuthUser{
		UserID: resp.GetUserId(),
		Email:  resp.GetEmail(),
		Role:   core.Role(resp.GetRole()),
	}, nil
}

func (c *Client) DeleteCredentials(ctx context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := c.client.DeleteCredentials(ctx, &authpb.DeleteCredentialsRequest{UserID: userID})
	return err
}
