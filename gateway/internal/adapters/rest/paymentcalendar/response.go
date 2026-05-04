package paymentcalendar

type PaymentResponse struct {
	ID          string `json:"id"`
	ProjectID   string `json:"project_id"`
	StageID     string `json:"stage_id,omitempty"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	PlannedDate string `json:"planned_date"`
	ActualDate  string `json:"actual_date,omitempty"`
	Description string `json:"description"`
	CreatedBy   string `json:"created_by,omitempty"`
	PaidBy      string `json:"paid_by,omitempty"`
	IsOverdue   bool   `json:"is_overdue"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

type ProjectFinancialSummaryResponse struct {
	ProjectID      string `json:"project_id"`
	PlannedIncome  string `json:"planned_income"`
	PlannedExpense string `json:"planned_expense"`
	PlannedBalance string `json:"planned_balance"`
	PaidIncome     string `json:"paid_income"`
	PaidExpense    string `json:"paid_expense"`
	PaidBalance    string `json:"paid_balance"`
	OverdueIncome  string `json:"overdue_income"`
	OverdueExpense string `json:"overdue_expense"`
	OverdueCount   int32  `json:"overdue_count"`
}

type ListPaymentsResponse struct {
	Payments []PaymentResponse `json:"payments"`
}
