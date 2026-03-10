package project

import (
	"context"
	"errors"
	"log/slog"
	"time"

	projectpb "github.com/Hubcher/project-management/contracts/gen/proto/project"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	log    *slog.Logger
	client projectpb.ProjectServiceClient
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
		client: projectpb.NewProjectServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.client.Ping(ctx, &emptypb.Empty{})
	return err
}

func (c *Client) CreateProject(ctx context.Context, req *projectpb.CreateProjectRequest) (*projectpb.Project, error) {
	if req == nil {
		return nil, errors.New("CreateProjectRequest cannot be nil")
	}

	resp, err := c.client.CreateProject(ctx, &projectpb.CreateProjectRequest{})
	if err != nil {
		c.log.Error("CreateProject failed", "error", err)
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetProject(ctx context.Context, req *projectpb.GetProjectRequest) (*projectpb.Project, error) {
	if req == nil {
		return nil, errors.New("get project request is nil")
	}

	resp, err := c.client.GetProject(ctx, req)
	if err != nil {
		c.log.Error("GetProject failed", "error", err)
		return nil, err
	}

	return resp, nil
}

func (c *Client) ListProject(ctx context.Context, req *projectpb.ListProjectsRequest) (*projectpb.ListProjectsResponse, error) {
	if req == nil {
		return nil, errors.New("list projects request is nil")
	}

	resp, err := c.client.ListProjects(ctx, req)
	if err != nil {
		c.log.Error("ListProjects failed", "error", err)
		return nil, err
	}

	return resp, nil
}

func (c *Client) UpdateProject(ctx context.Context, req *projectpb.UpdateProjectRequest) (*projectpb.Project, error) {
	if req == nil {
		return nil, errors.New("update project request is nil")
	}

	resp, err := c.client.UpdateProject(ctx, req)
	if err != nil {
		c.log.Error("UpdateProject failed", "error", err)
		return nil, err
	}

	return resp, nil
}

func (c *Client) DeleteProject(ctx context.Context, req *projectpb.DeleteProjectRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, errors.New("delete project request is nil")
	}

	resp, err := c.client.DeleteProject(ctx, req)
	if err != nil {
		c.log.Error("DeleteProject failed", "error", err)
		return nil, err
	}

	return resp, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
