package export

import (
	"context"
	"errors"
	"time"

	exportpb "github.com/Hubcher/project-management/contracts/gen/go/export"
	"github.com/Hubcher/project-management/gateway/internal/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client exportpb.ExportServiceClient
	conn   *grpc.ClientConn
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.CallContentSubtype(exportpb.JSONCodecName)),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.Config{BaseDelay: time.Second, Multiplier: 1.6, MaxDelay: 5 * time.Second},
			MinConnectTimeout: 10 * time.Second,
		}),
	)
	if err != nil {
		return nil, err
	}
	conn.Connect()
	return &Client{client: exportpb.NewExportServiceClient(conn), conn: conn}, nil
}

func (c *Client) Close() error { return c.conn.Close() }

func (c *Client) Ping(ctx context.Context) error {
	_, err := callWithTimeout(ctx, func(callCtx context.Context) (*exportpb.Empty, error) {
		return c.client.Ping(callCtx, &exportpb.Empty{})
	})
	return err
}

func (c *Client) BuildExport(ctx context.Context, input core.BuildExportInput) (*core.ExportedFile, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*exportpb.BuildExportResponse, error) {
		return c.client.BuildExport(callCtx, &exportpb.BuildExportRequest{
			ReportType:         input.ReportType,
			Format:             input.Format,
			ProjectId:          input.ProjectID,
			DateFrom:           input.DateFrom,
			DateTo:             input.DateTo,
			GroupBy:            input.GroupBy,
			PaymentType:        input.PaymentType,
			PaymentStatus:      input.PaymentStatus,
			OverdueOnly:        input.OverdueOnly,
			RequesterUserId:    input.RequesterUserID,
			IncludeAllProjects: input.IncludeAllProjects,
		})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoFile(resp)
}

func callWithTimeout[T any](ctx context.Context, fn func(context.Context) (T, error)) (T, error) {
	callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	return fn(callCtx)
}

func fromProtoFile(file *exportpb.BuildExportResponse) (*core.ExportedFile, error) {
	if file == nil {
		return nil, errors.New("export payload is empty")
	}
	return &core.ExportedFile{
		FileName:    file.FileName,
		ContentType: file.ContentType,
		Data:        file.Data,
	}, nil
}
