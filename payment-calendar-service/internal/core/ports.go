package core

import "context"

type PaymentRepository interface {
	CreatePayment(ctx context.Context, input CreatePaymentInput) (*Payment, error)
	GetPayment(ctx context.Context, id string) (*Payment, error)
	ListPayments(ctx context.Context, filter ListPaymentsFilter) ([]Payment, error)
	UpdatePayment(ctx context.Context, input UpdatePaymentInput) (*Payment, error)
	DeletePayment(ctx context.Context, id string) error
	MarkPaymentPaid(ctx context.Context, input MarkPaymentPaidInput) (*Payment, error)
}

type ProjectDirectory interface {
	GetProject(ctx context.Context, id string) (*ProjectRef, error)
	GetStage(ctx context.Context, id string) (*ProjectStageRef, error)
}

type PaymentCalendarService interface {
	CreatePayment(ctx context.Context, input CreatePaymentInput) (*Payment, error)
	GetPayment(ctx context.Context, id string) (*Payment, error)
	ListPayments(ctx context.Context, filter ListPaymentsFilter) ([]Payment, error)
	UpdatePayment(ctx context.Context, input UpdatePaymentInput) (*Payment, error)
	DeletePayment(ctx context.Context, id string) error
	MarkPaymentPaid(ctx context.Context, input MarkPaymentPaidInput) (*Payment, error)
	GetProjectSummary(ctx context.Context, projectID string) (*ProjectFinancialSummary, error)
}
