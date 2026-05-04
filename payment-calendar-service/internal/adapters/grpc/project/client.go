package project

import (
	"context"
	"time"

	projectpb "github.com/Hubcher/project-management/contracts/gen/go/project"
	"github.com/Hubcher/project-management/payment-calendar-service/internal/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Client struct {
	client projectpb.ProjectServiceClient
	conn   *grpc.ClientConn
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.CallContentSubtype(projectpb.JSONCodecName)),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.Config{BaseDelay: time.Second, Multiplier: 1.6, MaxDelay: 5 * time.Second},
			MinConnectTimeout: 10 * time.Second,
		}),
	)
	if err != nil {
		return nil, err
	}
	conn.Connect()
	return &Client{client: projectpb.NewProjectServiceClient(conn), conn: conn}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetProject(ctx context.Context, id string) (*core.ProjectRef, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*projectpb.Project, error) {
		return c.client.GetProject(callCtx, &projectpb.GetProjectRequest{Id: id})
	})
	if err != nil {
		return nil, mapProjectError(err, core.ErrProjectNotFound)
	}
	return &core.ProjectRef{ID: resp.ID}, nil
}

func (c *Client) GetStage(ctx context.Context, id string) (*core.ProjectStageRef, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*projectpb.ProjectStage, error) {
		return c.client.GetStage(callCtx, &projectpb.GetProjectStageRequest{Id: id})
	})
	if err != nil {
		return nil, mapProjectError(err, core.ErrStageNotFound)
	}
	return &core.ProjectStageRef{ID: resp.ID, ProjectID: resp.ProjectId}, nil
}

func callWithTimeout[T any](ctx context.Context, fn func(context.Context) (T, error)) (T, error) {
	callCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return fn(callCtx)
}

func mapProjectError(err error, notFound error) error {
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.NotFound:
			return notFound
		case codes.InvalidArgument:
			return core.ErrInvalidPayment
		}
	}
	return err
}
