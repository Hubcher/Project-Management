package paymentcalendar

import (
	"strings"

	"github.com/Hubcher/project-management/gateway/internal/core"
)

func toCreatePaymentInput(projectID, createdBy string, req CreatePaymentRequest) core.CreatePaymentInput {
	return core.CreatePaymentInput{
		ProjectID:   projectID,
		StageID:     req.StageID,
		Type:        req.Type,
		Amount:      req.Amount,
		Currency:    req.Currency,
		PlannedDate: req.PlannedDate,
		Description: req.Description,
		CreatedBy:   createdBy,
	}
}

func toListPaymentsInput(projectID string, query queryValues) core.ListPaymentsInput {
	return core.ListPaymentsInput{
		ProjectID:   projectID,
		StageID:     strings.TrimSpace(query.stageID),
		Type:        strings.TrimSpace(query.paymentType),
		Status:      strings.TrimSpace(query.status),
		DateFrom:    strings.TrimSpace(query.dateFrom),
		DateTo:      strings.TrimSpace(query.dateTo),
		OverdueOnly: query.overdueOnly,
	}
}

func toUpdatePaymentInput(id string, current *core.Payment, req UpdatePaymentRequest) core.UpdatePaymentInput {
	return core.UpdatePaymentInput{
		ID:          id,
		StageID:     mergeOptionalString(current.StageID, req.StageID),
		Type:        mergeOptionalString(current.Type, req.Type),
		Status:      mergeOptionalString(current.Status, req.Status),
		Amount:      mergeOptionalString(current.Amount, req.Amount),
		Currency:    mergeOptionalString(current.Currency, req.Currency),
		PlannedDate: mergeOptionalString(current.PlannedDate, req.PlannedDate),
		ActualDate:  mergeOptionalString(current.ActualDate, req.ActualDate),
		Description: mergeOptionalString(current.Description, req.Description),
		PaidBy:      current.PaidBy,
	}
}

func toPaymentResponse(payment *core.Payment) PaymentResponse {
	return PaymentResponse{
		ID:          payment.ID,
		ProjectID:   payment.ProjectID,
		StageID:     payment.StageID,
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
	}
}

func toListPaymentsResponse(payments []core.Payment) ListPaymentsResponse {
	resp := ListPaymentsResponse{Payments: make([]PaymentResponse, 0, len(payments))}
	for i := range payments {
		payment := payments[i]
		resp.Payments = append(resp.Payments, toPaymentResponse(&payment))
	}
	return resp
}

func toSummaryResponse(summary *core.ProjectFinancialSummary) ProjectFinancialSummaryResponse {
	return ProjectFinancialSummaryResponse{
		ProjectID:      summary.ProjectID,
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

func mergeOptionalString(current string, provided *string) string {
	if provided == nil {
		return current
	}
	return strings.TrimSpace(*provided)
}
