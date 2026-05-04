package paymentcalendar

import (
	"context"
	"errors"
	"time"

	paymentpb "github.com/Hubcher/project-management/contracts/gen/go/paymentcalendar"
	"github.com/Hubcher/project-management/gateway/internal/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client paymentpb.PaymentCalendarServiceClient
	conn   *grpc.ClientConn
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.CallContentSubtype(paymentpb.JSONCodecName)),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.Config{BaseDelay: time.Second, Multiplier: 1.6, MaxDelay: 5 * time.Second},
			MinConnectTimeout: 10 * time.Second,
		}),
	)
	if err != nil {
		return nil, err
	}
	conn.Connect()
	return &Client{client: paymentpb.NewPaymentCalendarServiceClient(conn), conn: conn}, nil
}

func (c *Client) Close() error { return c.conn.Close() }

func (c *Client) Ping(ctx context.Context) error {
	_, err := callWithTimeout(ctx, func(callCtx context.Context) (*paymentpb.Empty, error) {
		return c.client.Ping(callCtx, &paymentpb.Empty{})
	})
	return err
}

func (c *Client) CreatePayment(ctx context.Context, input core.CreatePaymentInput) (*core.Payment, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*paymentpb.Payment, error) {
		return c.client.CreatePayment(callCtx, &paymentpb.CreatePaymentRequest{
			ProjectId:   input.ProjectID,
			StageId:     input.StageID,
			Type:        input.Type,
			Amount:      input.Amount,
			Currency:    input.Currency,
			PlannedDate: input.PlannedDate,
			Description: input.Description,
			CreatedBy:   input.CreatedBy,
		})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoPayment(resp)
}

func (c *Client) GetPayment(ctx context.Context, id string) (*core.Payment, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*paymentpb.Payment, error) {
		return c.client.GetPayment(callCtx, &paymentpb.GetPaymentRequest{Id: id})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoPayment(resp)
}

func (c *Client) ListPayments(ctx context.Context, input core.ListPaymentsInput) ([]core.Payment, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*paymentpb.ListPaymentsResponse, error) {
		return c.client.ListPayments(callCtx, &paymentpb.ListPaymentsRequest{
			ProjectId:   input.ProjectID,
			StageId:     input.StageID,
			Type:        input.Type,
			Status:      input.Status,
			DateFrom:    input.DateFrom,
			DateTo:      input.DateTo,
			OverdueOnly: input.OverdueOnly,
		})
	})
	if err != nil {
		return nil, err
	}
	items := make([]core.Payment, 0, len(resp.Payments))
	for _, payment := range resp.Payments {
		item, convErr := fromProtoPayment(payment)
		if convErr != nil {
			return nil, convErr
		}
		items = append(items, *item)
	}
	return items, nil
}

func (c *Client) UpdatePayment(ctx context.Context, input core.UpdatePaymentInput) (*core.Payment, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*paymentpb.Payment, error) {
		return c.client.UpdatePayment(callCtx, &paymentpb.UpdatePaymentRequest{
			Id:          input.ID,
			StageId:     input.StageID,
			Type:        input.Type,
			Status:      input.Status,
			Amount:      input.Amount,
			Currency:    input.Currency,
			PlannedDate: input.PlannedDate,
			ActualDate:  input.ActualDate,
			Description: input.Description,
			PaidBy:      input.PaidBy,
		})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoPayment(resp)
}

func (c *Client) DeletePayment(ctx context.Context, id string) error {
	_, err := callWithTimeout(ctx, func(callCtx context.Context) (*paymentpb.Empty, error) {
		return c.client.DeletePayment(callCtx, &paymentpb.DeletePaymentRequest{Id: id})
	})
	return err
}

func (c *Client) MarkPaymentPaid(ctx context.Context, input core.MarkPaymentPaidInput) (*core.Payment, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*paymentpb.Payment, error) {
		return c.client.MarkPaymentPaid(callCtx, &paymentpb.MarkPaymentPaidRequest{
			Id:         input.ID,
			ActualDate: input.ActualDate,
			PaidBy:     input.PaidBy,
		})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoPayment(resp)
}

func (c *Client) GetProjectSummary(ctx context.Context, projectID string) (*core.ProjectFinancialSummary, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*paymentpb.ProjectFinancialSummary, error) {
		return c.client.GetProjectSummary(callCtx, &paymentpb.GetProjectSummaryRequest{ProjectId: projectID})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoSummary(resp)
}

func callWithTimeout[T any](ctx context.Context, fn func(context.Context) (T, error)) (T, error) {
	callCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return fn(callCtx)
}

func fromProtoPayment(payment *paymentpb.Payment) (*core.Payment, error) {
	if payment == nil {
		return nil, errors.New("payment payload is empty")
	}
	return &core.Payment{
		ID:          payment.ID,
		ProjectID:   payment.ProjectId,
		StageID:     payment.StageId,
		Type:        payment.Type,
		Status:      payment.Status,
		Amount:      payment.Amount,
		Currency:    payment.Currency,
		PlannedDate: payment.PlannedDate,
		ActualDate:  payment.ActualDate,
		Description: payment.Description,
		CreatedBy:   payment.CreatedBy,
		PaidBy:      payment.PaidBy,
		IsOverdue:   payment.IsOverdue,
		CreatedAt:   payment.CreatedAt,
		UpdatedAt:   payment.UpdatedAt,
	}, nil
}

func fromProtoSummary(summary *paymentpb.ProjectFinancialSummary) (*core.ProjectFinancialSummary, error) {
	if summary == nil {
		return nil, errors.New("payment summary payload is empty")
	}
	return &core.ProjectFinancialSummary{
		ProjectID:      summary.ProjectId,
		PlannedIncome:  summary.PlannedIncome,
		PlannedExpense: summary.PlannedExpense,
		PlannedBalance: summary.PlannedBalance,
		PaidIncome:     summary.PaidIncome,
		PaidExpense:    summary.PaidExpense,
		PaidBalance:    summary.PaidBalance,
		OverdueIncome:  summary.OverdueIncome,
		OverdueExpense: summary.OverdueExpense,
		OverdueCount:   summary.OverdueCount,
	}, nil
}
