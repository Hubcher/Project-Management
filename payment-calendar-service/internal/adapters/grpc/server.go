package grpc

import (
	"context"
	"errors"
	"time"

	paymentpb "github.com/Hubcher/project-management/contracts/gen/go/paymentcalendar"
	"github.com/Hubcher/project-management/payment-calendar-service/internal/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	paymentpb.UnimplementedPaymentCalendarServiceServer
	service core.PaymentCalendarService
}

func NewServer(service core.PaymentCalendarService) *Server {
	return &Server{service: service}
}

func (s *Server) Ping(_ context.Context, _ *paymentpb.Empty) (*paymentpb.Empty, error) {
	return &paymentpb.Empty{}, nil
}

func (s *Server) CreatePayment(ctx context.Context, req *paymentpb.CreatePaymentRequest) (*paymentpb.Payment, error) {
	plannedDate, err := parseDate(req.PlannedDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid planned_date")
	}
	payment, err := s.service.CreatePayment(ctx, core.CreatePaymentInput{
		ProjectID:   req.ProjectId,
		StageID:     req.StageId,
		Type:        core.PaymentType(req.Type),
		Amount:      req.Amount,
		Currency:    req.Currency,
		PlannedDate: plannedDate,
		Description: req.Description,
		CreatedBy:   req.CreatedBy,
	})
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoPayment(payment), nil
}

func (s *Server) GetPayment(ctx context.Context, req *paymentpb.GetPaymentRequest) (*paymentpb.Payment, error) {
	payment, err := s.service.GetPayment(ctx, req.Id)
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoPayment(payment), nil
}

func (s *Server) ListPayments(ctx context.Context, req *paymentpb.ListPaymentsRequest) (*paymentpb.ListPaymentsResponse, error) {
	dateFrom, err := parseDate(req.DateFrom)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid date_from")
	}
	dateTo, err := parseDate(req.DateTo)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid date_to")
	}

	payments, err := s.service.ListPayments(ctx, core.ListPaymentsFilter{
		ProjectID:   req.ProjectId,
		StageID:     req.StageId,
		Type:        core.PaymentType(req.Type),
		Status:      core.PaymentStatus(req.Status),
		DateFrom:    dateFrom,
		DateTo:      dateTo,
		OverdueOnly: req.OverdueOnly,
	})
	if err != nil {
		return nil, mapCoreError(err)
	}
	resp := &paymentpb.ListPaymentsResponse{Payments: make([]*paymentpb.Payment, 0, len(payments))}
	for i := range payments {
		resp.Payments = append(resp.Payments, toProtoPayment(&payments[i]))
	}
	return resp, nil
}

func (s *Server) UpdatePayment(ctx context.Context, req *paymentpb.UpdatePaymentRequest) (*paymentpb.Payment, error) {
	plannedDate, err := parseDate(req.PlannedDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid planned_date")
	}
	actualDate, err := parseDate(req.ActualDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid actual_date")
	}

	payment, err := s.service.UpdatePayment(ctx, core.UpdatePaymentInput{
		ID:          req.Id,
		StageID:     req.StageId,
		Type:        core.PaymentType(req.Type),
		Status:      core.PaymentStatus(req.Status),
		Amount:      req.Amount,
		Currency:    req.Currency,
		PlannedDate: plannedDate,
		ActualDate:  actualDate,
		Description: req.Description,
		PaidBy:      req.PaidBy,
	})
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoPayment(payment), nil
}

func (s *Server) DeletePayment(ctx context.Context, req *paymentpb.DeletePaymentRequest) (*paymentpb.Empty, error) {
	if err := s.service.DeletePayment(ctx, req.Id); err != nil {
		return nil, mapCoreError(err)
	}
	return &paymentpb.Empty{}, nil
}

func (s *Server) MarkPaymentPaid(ctx context.Context, req *paymentpb.MarkPaymentPaidRequest) (*paymentpb.Payment, error) {
	actualDate, err := parseDate(req.ActualDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid actual_date")
	}
	payment, err := s.service.MarkPaymentPaid(ctx, core.MarkPaymentPaidInput{
		ID:         req.Id,
		ActualDate: actualDate,
		PaidBy:     req.PaidBy,
	})
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoPayment(payment), nil
}

func (s *Server) GetProjectSummary(ctx context.Context, req *paymentpb.GetProjectSummaryRequest) (*paymentpb.ProjectFinancialSummary, error) {
	summary, err := s.service.GetProjectSummary(ctx, req.ProjectId)
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoSummary(summary), nil
}

func parseDate(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, err
	}
	day := parsed.UTC()
	normalized := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	return &normalized, nil
}

func toProtoPayment(payment *core.Payment) *paymentpb.Payment {
	return &paymentpb.Payment{
		ID:          payment.ID,
		ProjectId:   payment.ProjectID,
		StageId:     payment.StageID,
		Type:        string(payment.Type),
		Status:      string(payment.Status),
		Amount:      payment.Amount,
		Currency:    payment.Currency,
		PlannedDate: formatDate(&payment.PlannedDate),
		ActualDate:  formatDate(payment.ActualDate),
		Description: payment.Description,
		CreatedBy:   payment.CreatedBy,
		PaidBy:      payment.PaidBy,
		IsOverdue:   payment.IsOverdue,
		CreatedAt:   formatTimestamp(&payment.CreatedAt),
		UpdatedAt:   formatTimestamp(&payment.UpdatedAt),
	}
}

func toProtoSummary(summary *core.ProjectFinancialSummary) *paymentpb.ProjectFinancialSummary {
	return &paymentpb.ProjectFinancialSummary{
		ProjectId:      summary.ProjectID,
		PlannedIncome:  summary.PlannedIncome,
		PlannedExpense: summary.PlannedExpense,
		PlannedBalance: summary.PlannedBalance,
		PaidIncome:     summary.PaidIncome,
		PaidExpense:    summary.PaidExpense,
		PaidBalance:    summary.PaidBalance,
		OverdueIncome:  summary.OverdueIncome,
		OverdueExpense: summary.OverdueExpense,
		OverdueCount:   summary.OverdueCount,
	}
}

func formatDate(value *time.Time) string {
	if value == nil || value.IsZero() {
		return ""
	}
	return value.UTC().Format("2006-01-02")
}

func formatTimestamp(value *time.Time) string {
	if value == nil || value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339)
}

func mapCoreError(err error) error {
	switch {
	case errors.Is(err, core.ErrInvalidPayment):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, core.ErrPaymentNotFound), errors.Is(err, core.ErrProjectNotFound), errors.Is(err, core.ErrStageNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, core.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
