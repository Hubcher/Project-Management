package core

import "time"

type PaymentType string

type PaymentStatus string

const (
	PaymentTypeIncome  PaymentType = "income"
	PaymentTypeExpense PaymentType = "expense"

	PaymentStatusPlanned   PaymentStatus = "planned"
	PaymentStatusPaid      PaymentStatus = "paid"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

type Payment struct {
	ID          string
	ProjectID   string
	StageID     string
	Type        PaymentType
	Status      PaymentStatus
	Amount      string
	Currency    string
	PlannedDate time.Time
	ActualDate  *time.Time
	Description string
	CreatedBy   string
	PaidBy      string
	IsOverdue   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreatePaymentInput struct {
	ProjectID   string
	StageID     string
	Type        PaymentType
	Amount      string
	Currency    string
	PlannedDate *time.Time
	Description string
	CreatedBy   string
}

type UpdatePaymentInput struct {
	ID          string
	StageID     string
	Type        PaymentType
	Status      PaymentStatus
	Amount      string
	Currency    string
	PlannedDate *time.Time
	ActualDate  *time.Time
	Description string
	PaidBy      string
}

type MarkPaymentPaidInput struct {
	ID         string
	ActualDate *time.Time
	PaidBy     string
}

type ListPaymentsFilter struct {
	ProjectID   string
	StageID     string
	Type        PaymentType
	Status      PaymentStatus
	DateFrom    *time.Time
	DateTo      *time.Time
	OverdueOnly bool
}

type ProjectFinancialSummary struct {
	ProjectID       string
	PlannedIncome   string
	PlannedExpense  string
	PlannedBalance  string
	PaidIncome      string
	PaidExpense     string
	PaidBalance     string
	OverdueIncome   string
	OverdueExpense  string
	OverdueCount    int32
}

type ProjectRef struct {
	ID string
}

type ProjectStageRef struct {
	ID        string
	ProjectID string
}
