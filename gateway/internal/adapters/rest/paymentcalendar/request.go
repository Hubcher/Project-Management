package paymentcalendar

type CreatePaymentRequest struct {
	StageID     string `json:"stage_id"`
	Type        string `json:"type"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	PlannedDate string `json:"planned_date"`
	Description string `json:"description"`
}

type UpdatePaymentRequest struct {
	StageID     *string `json:"stage_id"`
	Type        *string `json:"type"`
	Status      *string `json:"status"`
	Amount      *string `json:"amount"`
	Currency    *string `json:"currency"`
	PlannedDate *string `json:"planned_date"`
	ActualDate  *string `json:"actual_date"`
	Description *string `json:"description"`
}

type MarkPaymentPaidRequest struct {
	ActualDate string `json:"actual_date"`
}
