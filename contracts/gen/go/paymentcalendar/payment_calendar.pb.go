package paymentcalendarpb

import (
	"encoding/json"

	"google.golang.org/grpc/encoding"
)

const JSONCodecName = "json"

type jsonCodec struct{}

func (jsonCodec) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (jsonCodec) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (jsonCodec) Name() string {
	return JSONCodecName
}

func init() {
	encoding.RegisterCodec(jsonCodec{})
}

type Empty struct{}

type Payment struct {
	ID          string `json:"id,omitempty"`
	ProjectId   string `json:"project_id,omitempty"`
	StageId     string `json:"stage_id,omitempty"`
	Type        string `json:"type,omitempty"`
	Status      string `json:"status,omitempty"`
	Amount      string `json:"amount,omitempty"`
	Currency    string `json:"currency,omitempty"`
	PlannedDate string `json:"planned_date,omitempty"`
	ActualDate  string `json:"actual_date,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedBy   string `json:"created_by,omitempty"`
	PaidBy      string `json:"paid_by,omitempty"`
	IsOverdue   bool   `json:"is_overdue,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

type CreatePaymentRequest struct {
	ProjectId   string `json:"project_id,omitempty"`
	StageId     string `json:"stage_id,omitempty"`
	Type        string `json:"type,omitempty"`
	Amount      string `json:"amount,omitempty"`
	Currency    string `json:"currency,omitempty"`
	PlannedDate string `json:"planned_date,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedBy   string `json:"created_by,omitempty"`
}

type GetPaymentRequest struct {
	Id string `json:"id,omitempty"`
}

type ListPaymentsRequest struct {
	ProjectId   string `json:"project_id,omitempty"`
	StageId     string `json:"stage_id,omitempty"`
	Type        string `json:"type,omitempty"`
	Status      string `json:"status,omitempty"`
	DateFrom    string `json:"date_from,omitempty"`
	DateTo      string `json:"date_to,omitempty"`
	OverdueOnly bool   `json:"overdue_only,omitempty"`
}

type ListPaymentsResponse struct {
	Payments []*Payment `json:"payments,omitempty"`
}

type UpdatePaymentRequest struct {
	Id          string `json:"id,omitempty"`
	StageId     string `json:"stage_id,omitempty"`
	Type        string `json:"type,omitempty"`
	Status      string `json:"status,omitempty"`
	Amount      string `json:"amount,omitempty"`
	Currency    string `json:"currency,omitempty"`
	PlannedDate string `json:"planned_date,omitempty"`
	ActualDate  string `json:"actual_date,omitempty"`
	Description string `json:"description,omitempty"`
	PaidBy      string `json:"paid_by,omitempty"`
}

type DeletePaymentRequest struct {
	Id string `json:"id,omitempty"`
}

type MarkPaymentPaidRequest struct {
	Id         string `json:"id,omitempty"`
	ActualDate string `json:"actual_date,omitempty"`
	PaidBy     string `json:"paid_by,omitempty"`
}

type GetProjectSummaryRequest struct {
	ProjectId string `json:"project_id,omitempty"`
}

type ProjectFinancialSummary struct {
	ProjectId       string `json:"project_id,omitempty"`
	PlannedIncome   string `json:"planned_income,omitempty"`
	PlannedExpense  string `json:"planned_expense,omitempty"`
	PlannedBalance  string `json:"planned_balance,omitempty"`
	PaidIncome      string `json:"paid_income,omitempty"`
	PaidExpense     string `json:"paid_expense,omitempty"`
	PaidBalance     string `json:"paid_balance,omitempty"`
	OverdueIncome   string `json:"overdue_income,omitempty"`
	OverdueExpense  string `json:"overdue_expense,omitempty"`
	OverdueCount    int32  `json:"overdue_count,omitempty"`
}
