package report

import (
	"context"
	"errors"
	"time"

	reportpb "github.com/Hubcher/project-management/contracts/gen/go/report"
	"github.com/Hubcher/project-management/gateway/internal/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client reportpb.ReportServiceClient
	conn   *grpc.ClientConn
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.CallContentSubtype(reportpb.JSONCodecName)),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.Config{BaseDelay: time.Second, Multiplier: 1.6, MaxDelay: 5 * time.Second},
			MinConnectTimeout: 10 * time.Second,
		}),
	)
	if err != nil {
		return nil, err
	}
	conn.Connect()
	return &Client{client: reportpb.NewReportServiceClient(conn), conn: conn}, nil
}

func (c *Client) Close() error { return c.conn.Close() }

func (c *Client) Ping(ctx context.Context) error {
	_, err := callWithTimeout(ctx, func(callCtx context.Context) (*reportpb.Empty, error) {
		return c.client.Ping(callCtx, &reportpb.Empty{})
	})
	return err
}

func (c *Client) CreateReport(ctx context.Context, input core.CreateDailyReportInput) (*core.DailyReport, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*reportpb.DailyReport, error) {
		return c.client.CreateReport(callCtx, &reportpb.CreateDailyReportRequest{
			UserId:     input.UserID,
			ReportDate: input.ReportDate,
			Status:     input.Status,
			Summary:    input.Summary,
		})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoReport(resp)
}

func (c *Client) GetReport(ctx context.Context, id string) (*core.DailyReport, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*reportpb.DailyReport, error) {
		return c.client.GetReport(callCtx, &reportpb.GetDailyReportRequest{Id: id})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoReport(resp)
}

func (c *Client) ListReports(ctx context.Context, input core.ListDailyReportsInput) ([]core.DailyReport, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*reportpb.ListDailyReportsResponse, error) {
		return c.client.ListReports(callCtx, &reportpb.ListDailyReportsRequest{
			UserId:   input.UserID,
			Status:   input.Status,
			DateFrom: input.DateFrom,
			DateTo:   input.DateTo,
		})
	})
	if err != nil {
		return nil, err
	}
	items := make([]core.DailyReport, 0, len(resp.Reports))
	for _, report := range resp.Reports {
		item, convErr := fromProtoReport(report)
		if convErr != nil {
			return nil, convErr
		}
		items = append(items, *item)
	}
	return items, nil
}

func (c *Client) UpdateReport(ctx context.Context, input core.UpdateDailyReportInput) (*core.DailyReport, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*reportpb.DailyReport, error) {
		return c.client.UpdateReport(callCtx, &reportpb.UpdateDailyReportRequest{
			Id:      input.ID,
			Status:  input.Status,
			Summary: input.Summary,
		})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoReport(resp)
}

func (c *Client) DeleteReport(ctx context.Context, id string) error {
	_, err := callWithTimeout(ctx, func(callCtx context.Context) (*reportpb.Empty, error) {
		return c.client.DeleteReport(callCtx, &reportpb.DeleteDailyReportRequest{Id: id})
	})
	return err
}

func (c *Client) CreateEntry(ctx context.Context, input core.CreateDailyReportEntryInput) (*core.DailyReportEntry, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*reportpb.DailyReportEntry, error) {
		return c.client.CreateEntry(callCtx, &reportpb.CreateDailyReportEntryRequest{
			ReportId:    input.ReportID,
			ProjectId:   input.ProjectID,
			StageId:     input.StageID,
			WorkType:    input.WorkType,
			Description: input.Description,
			HoursSpent:  input.HoursSpent,
		})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoEntry(resp)
}

func (c *Client) GetEntry(ctx context.Context, id string) (*core.DailyReportEntry, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*reportpb.DailyReportEntry, error) {
		return c.client.GetEntry(callCtx, &reportpb.GetDailyReportEntryRequest{Id: id})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoEntry(resp)
}

func (c *Client) ListEntries(ctx context.Context, reportID string) ([]core.DailyReportEntry, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*reportpb.ListDailyReportEntriesResponse, error) {
		return c.client.ListEntries(callCtx, &reportpb.ListDailyReportEntriesRequest{ReportId: reportID})
	})
	if err != nil {
		return nil, err
	}
	items := make([]core.DailyReportEntry, 0, len(resp.Entries))
	for _, entry := range resp.Entries {
		item, convErr := fromProtoEntry(entry)
		if convErr != nil {
			return nil, convErr
		}
		items = append(items, *item)
	}
	return items, nil
}

func (c *Client) UpdateEntry(ctx context.Context, input core.UpdateDailyReportEntryInput) (*core.DailyReportEntry, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*reportpb.DailyReportEntry, error) {
		return c.client.UpdateEntry(callCtx, &reportpb.UpdateDailyReportEntryRequest{
			Id:          input.ID,
			ProjectId:   input.ProjectID,
			StageId:     input.StageID,
			WorkType:    input.WorkType,
			Description: input.Description,
			HoursSpent:  input.HoursSpent,
		})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoEntry(resp)
}

func (c *Client) DeleteEntry(ctx context.Context, id string) error {
	_, err := callWithTimeout(ctx, func(callCtx context.Context) (*reportpb.Empty, error) {
		return c.client.DeleteEntry(callCtx, &reportpb.DeleteDailyReportEntryRequest{Id: id})
	})
	return err
}

func (c *Client) CreateComment(ctx context.Context, input core.CreateDailyReportCommentInput) (*core.DailyReportComment, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*reportpb.DailyReportComment, error) {
		return c.client.CreateComment(callCtx, &reportpb.CreateDailyReportCommentRequest{
			ReportId:     input.ReportID,
			AuthorUserId: input.AuthorUserID,
			Comment:      input.Comment,
		})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoComment(resp)
}

func (c *Client) GetComment(ctx context.Context, id string) (*core.DailyReportComment, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*reportpb.DailyReportComment, error) {
		return c.client.GetComment(callCtx, &reportpb.GetDailyReportCommentRequest{Id: id})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoComment(resp)
}

func (c *Client) ListComments(ctx context.Context, reportID string) ([]core.DailyReportComment, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*reportpb.ListDailyReportCommentsResponse, error) {
		return c.client.ListComments(callCtx, &reportpb.ListDailyReportCommentsRequest{ReportId: reportID})
	})
	if err != nil {
		return nil, err
	}
	items := make([]core.DailyReportComment, 0, len(resp.Comments))
	for _, comment := range resp.Comments {
		item, convErr := fromProtoComment(comment)
		if convErr != nil {
			return nil, convErr
		}
		items = append(items, *item)
	}
	return items, nil
}

func (c *Client) UpdateComment(ctx context.Context, input core.UpdateDailyReportCommentInput) (*core.DailyReportComment, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*reportpb.DailyReportComment, error) {
		return c.client.UpdateComment(callCtx, &reportpb.UpdateDailyReportCommentRequest{
			Id:      input.ID,
			Comment: input.Comment,
		})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoComment(resp)
}

func (c *Client) DeleteComment(ctx context.Context, id string) error {
	_, err := callWithTimeout(ctx, func(callCtx context.Context) (*reportpb.Empty, error) {
		return c.client.DeleteComment(callCtx, &reportpb.DeleteDailyReportCommentRequest{Id: id})
	})
	return err
}

func callWithTimeout[T any](ctx context.Context, fn func(context.Context) (T, error)) (T, error) {
	callCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return fn(callCtx)
}

func fromProtoReport(report *reportpb.DailyReport) (*core.DailyReport, error) {
	if report == nil {
		return nil, errors.New("daily report payload is empty")
	}
	return &core.DailyReport{ID: report.ID, UserID: report.UserId, ReportDate: report.ReportDate, Status: report.Status, TotalHours: report.TotalHours, Summary: report.Summary, CreatedAt: report.CreatedAt, UpdatedAt: report.UpdatedAt}, nil
}

func fromProtoEntry(entry *reportpb.DailyReportEntry) (*core.DailyReportEntry, error) {
	if entry == nil {
		return nil, errors.New("daily report entry payload is empty")
	}
	return &core.DailyReportEntry{ID: entry.ID, ReportID: entry.ReportId, ProjectID: entry.ProjectId, StageID: entry.StageId, WorkType: entry.WorkType, Description: entry.Description, HoursSpent: entry.HoursSpent, CreatedAt: entry.CreatedAt, UpdatedAt: entry.UpdatedAt}, nil
}

func fromProtoComment(comment *reportpb.DailyReportComment) (*core.DailyReportComment, error) {
	if comment == nil {
		return nil, errors.New("daily report comment payload is empty")
	}
	return &core.DailyReportComment{ID: comment.ID, ReportID: comment.ReportId, AuthorUserID: comment.AuthorUserId, Comment: comment.Comment, CreatedAt: comment.CreatedAt}, nil
}