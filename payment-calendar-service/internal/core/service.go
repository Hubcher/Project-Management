package core

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"
)

type service struct {
	repo     PaymentRepository
	projects ProjectDirectory
	clock    func() time.Time
}

func NewService(repo PaymentRepository, projects ProjectDirectory) PaymentCalendarService {
	return &service{
		repo:     repo,
		projects: projects,
		clock:    func() time.Time { return time.Now().UTC() },
	}
}

func (s *service) CreatePayment(ctx context.Context, input CreatePaymentInput) (*Payment, error) {
	input = normalizeCreatePaymentInput(input)
	if err := validateCreatePaymentInput(input); err != nil {
		return nil, err
	}
	if err := s.ensureKnownProjectAndStage(ctx, input.ProjectID, input.StageID); err != nil {
		return nil, err
	}

	payment, err := s.repo.CreatePayment(ctx, input)
	if err != nil {
		return nil, err
	}
	return s.withComputedFields(payment), nil
}

func (s *service) GetPayment(ctx context.Context, id string) (*Payment, error) {
	id = strings.TrimSpace(id)
	if !isValidUUID(id) {
		return nil, ErrInvalidPayment
	}
	payment, err := s.repo.GetPayment(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.withComputedFields(payment), nil
}

func (s *service) ListPayments(ctx context.Context, filter ListPaymentsFilter) ([]Payment, error) {
	filter = normalizeListPaymentsFilter(filter)
	if err := validateListPaymentsFilter(filter); err != nil {
		return nil, err
	}
	if err := s.ensureKnownProjectAndStage(ctx, filter.ProjectID, filter.StageID); err != nil {
		return nil, err
	}

	payments, err := s.repo.ListPayments(ctx, filter)
	if err != nil {
		return nil, err
	}
	result := make([]Payment, 0, len(payments))
	for i := range payments {
		payment := *s.withComputedFields(&payments[i])
		if filter.OverdueOnly && !payment.IsOverdue {
			continue
		}
		result = append(result, payment)
	}
	return result, nil
}

func (s *service) UpdatePayment(ctx context.Context, input UpdatePaymentInput) (*Payment, error) {
	input = normalizeUpdatePaymentInput(input)
	if !isValidUUID(input.ID) {
		return nil, ErrInvalidPayment
	}

	current, err := s.repo.GetPayment(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if err = validateUpdatePaymentInput(input); err != nil {
		return nil, err
	}
	if err = s.ensureKnownProjectAndStage(ctx, current.ProjectID, input.StageID); err != nil {
		return nil, err
	}

	payment, err := s.repo.UpdatePayment(ctx, input)
	if err != nil {
		return nil, err
	}
	return s.withComputedFields(payment), nil
}

func (s *service) DeletePayment(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if !isValidUUID(id) {
		return ErrInvalidPayment
	}
	return s.repo.DeletePayment(ctx, id)
}

func (s *service) MarkPaymentPaid(ctx context.Context, input MarkPaymentPaidInput) (*Payment, error) {
	input = normalizeMarkPaymentPaidInput(input, s.clock())
	if !isValidUUID(input.ID) || input.ActualDate == nil {
		return nil, ErrInvalidPayment
	}
	if input.PaidBy != "" && !isValidUUID(input.PaidBy) {
		return nil, ErrInvalidPayment
	}

	current, err := s.repo.GetPayment(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if current.Status == PaymentStatusCancelled {
		return nil, ErrInvalidPayment
	}

	payment, err := s.repo.MarkPaymentPaid(ctx, input)
	if err != nil {
		return nil, err
	}
	return s.withComputedFields(payment), nil
}

func (s *service) GetProjectSummary(ctx context.Context, projectID string) (*ProjectFinancialSummary, error) {
	projectID = strings.TrimSpace(projectID)
	if !isValidUUID(projectID) {
		return nil, ErrInvalidPayment
	}
	if err := s.ensureKnownProjectAndStage(ctx, projectID, ""); err != nil {
		return nil, err
	}

	payments, err := s.ListPayments(ctx, ListPaymentsFilter{ProjectID: projectID})
	if err != nil {
		return nil, err
	}
	return buildProjectSummary(projectID, payments), nil
}

func (s *service) ensureKnownProjectAndStage(ctx context.Context, projectID, stageID string) error {
	if s.projects == nil {
		return nil
	}
	if _, err := s.projects.GetProject(ctx, projectID); err != nil {
		return err
	}
	stageID = strings.TrimSpace(stageID)
	if stageID == "" {
		return nil
	}
	stage, err := s.projects.GetStage(ctx, stageID)
	if err != nil {
		return err
	}
	if stage.ProjectID != projectID {
		return ErrInvalidPayment
	}
	return nil
}

func (s *service) withComputedFields(payment *Payment) *Payment {
	if payment == nil {
		return nil
	}
	payment.IsOverdue = isOverdue(payment.Status, payment.PlannedDate, s.clock())
	return payment
}

func validateCreatePaymentInput(input CreatePaymentInput) error {
	if !isValidUUID(input.ProjectID) || input.PlannedDate == nil {
		return ErrInvalidPayment
	}
	if input.StageID != "" && !isValidUUID(input.StageID) {
		return ErrInvalidPayment
	}
	if input.CreatedBy != "" && !isValidUUID(input.CreatedBy) {
		return ErrInvalidPayment
	}
	if !isValidPaymentType(input.Type) || !isPositiveMoney(input.Amount) || !isValidCurrency(input.Currency) {
		return ErrInvalidPayment
	}
	return nil
}

func validateUpdatePaymentInput(input UpdatePaymentInput) error {
	if input.StageID != "" && !isValidUUID(input.StageID) {
		return ErrInvalidPayment
	}
	if input.PaidBy != "" && !isValidUUID(input.PaidBy) {
		return ErrInvalidPayment
	}
	if input.PlannedDate == nil || !isValidPaymentType(input.Type) || !isValidPaymentStatus(input.Status) {
		return ErrInvalidPayment
	}
	if !isPositiveMoney(input.Amount) || !isValidCurrency(input.Currency) {
		return ErrInvalidPayment
	}
	if input.Status == PaymentStatusPaid && input.ActualDate == nil {
		return ErrInvalidPayment
	}
	return nil
}

func validateListPaymentsFilter(filter ListPaymentsFilter) error {
	if !isValidUUID(filter.ProjectID) {
		return ErrInvalidPayment
	}
	if filter.StageID != "" && !isValidUUID(filter.StageID) {
		return ErrInvalidPayment
	}
	if filter.Type != "" && !isValidPaymentType(filter.Type) {
		return ErrInvalidPayment
	}
	if filter.Status != "" && !isValidPaymentStatus(filter.Status) {
		return ErrInvalidPayment
	}
	if !isValidDateRange(filter.DateFrom, filter.DateTo) {
		return ErrInvalidPayment
	}
	return nil
}

func normalizeCreatePaymentInput(input CreatePaymentInput) CreatePaymentInput {
	input.ProjectID = strings.TrimSpace(input.ProjectID)
	input.StageID = strings.TrimSpace(input.StageID)
	input.Type = normalizePaymentType(input.Type)
	input.Amount = normalizeMoney(input.Amount)
	input.Currency = normalizeCurrency(input.Currency)
	input.PlannedDate = normalizeDate(input.PlannedDate)
	input.Description = strings.TrimSpace(input.Description)
	input.CreatedBy = strings.TrimSpace(input.CreatedBy)
	return input
}

func normalizeUpdatePaymentInput(input UpdatePaymentInput) UpdatePaymentInput {
	input.ID = strings.TrimSpace(input.ID)
	input.StageID = strings.TrimSpace(input.StageID)
	input.Type = normalizePaymentType(input.Type)
	input.Status = normalizePaymentStatus(input.Status)
	input.Amount = normalizeMoney(input.Amount)
	input.Currency = normalizeCurrency(input.Currency)
	input.PlannedDate = normalizeDate(input.PlannedDate)
	input.ActualDate = normalizeDate(input.ActualDate)
	input.Description = strings.TrimSpace(input.Description)
	input.PaidBy = strings.TrimSpace(input.PaidBy)
	if input.Status != PaymentStatusPaid {
		input.ActualDate = nil
		input.PaidBy = ""
	}
	return input
}

func normalizeMarkPaymentPaidInput(input MarkPaymentPaidInput, now time.Time) MarkPaymentPaidInput {
	input.ID = strings.TrimSpace(input.ID)
	input.ActualDate = normalizeDate(input.ActualDate)
	input.PaidBy = strings.TrimSpace(input.PaidBy)
	if input.ActualDate == nil {
		input.ActualDate = normalizeDate(&now)
	}
	return input
}

func normalizeListPaymentsFilter(filter ListPaymentsFilter) ListPaymentsFilter {
	filter.ProjectID = strings.TrimSpace(filter.ProjectID)
	filter.StageID = strings.TrimSpace(filter.StageID)
	filter.Type = normalizePaymentType(filter.Type)
	filter.Status = normalizePaymentStatus(filter.Status)
	filter.DateFrom = normalizeDate(filter.DateFrom)
	filter.DateTo = normalizeDate(filter.DateTo)
	return filter
}

func normalizePaymentType(value PaymentType) PaymentType {
	return PaymentType(strings.ToLower(strings.TrimSpace(string(value))))
}

func normalizePaymentStatus(value PaymentStatus) PaymentStatus {
	return PaymentStatus(strings.ToLower(strings.TrimSpace(string(value))))
}

func normalizeMoney(value string) string {
	return strings.TrimSpace(value)
}

func normalizeCurrency(value string) string {
	value = strings.ToUpper(strings.TrimSpace(value))
	if value == "" {
		return "RUB"
	}
	return value
}

func normalizeDate(value *time.Time) *time.Time {
	if value == nil || value.IsZero() {
		return nil
	}
	day := value.UTC()
	normalized := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	return &normalized
}

func isValidUUID(value string) bool {
	_, err := uuid.Parse(strings.TrimSpace(value))
	return err == nil
}

func isValidPaymentType(value PaymentType) bool {
	switch value {
	case PaymentTypeIncome, PaymentTypeExpense:
		return true
	default:
		return false
	}
}

func isValidPaymentStatus(value PaymentStatus) bool {
	switch value {
	case PaymentStatusPlanned, PaymentStatusPaid, PaymentStatusCancelled:
		return true
	default:
		return false
	}
}

func isPositiveMoney(value string) bool {
	var amount big.Rat
	if _, ok := amount.SetString(strings.TrimSpace(value)); !ok {
		return false
	}
	return amount.Sign() > 0
}

func isValidCurrency(value string) bool {
	value = strings.TrimSpace(value)
	return len(value) == 3
}

func isValidDateRange(start, end *time.Time) bool {
	if start == nil || end == nil {
		return true
	}
	return !end.Before(*start)
}

func isOverdue(status PaymentStatus, plannedDate time.Time, now time.Time) bool {
	if status != PaymentStatusPlanned || plannedDate.IsZero() {
		return false
	}
	today := normalizeDate(&now)
	if today == nil {
		return false
	}
	return plannedDate.Before(*today)
}

func buildProjectSummary(projectID string, payments []Payment) *ProjectFinancialSummary {
	var plannedIncome, plannedExpense big.Rat
	var paidIncome, paidExpense big.Rat
	var overdueIncome, overdueExpense big.Rat
	var overdueCount int32

	for i := range payments {
		payment := payments[i]
		if payment.Status == PaymentStatusCancelled {
			continue
		}

		switch payment.Type {
		case PaymentTypeIncome:
			addMoney(&plannedIncome, payment.Amount)
			if payment.Status == PaymentStatusPaid {
				addMoney(&paidIncome, payment.Amount)
			}
			if payment.IsOverdue {
				addMoney(&overdueIncome, payment.Amount)
				overdueCount++
			}
		case PaymentTypeExpense:
			addMoney(&plannedExpense, payment.Amount)
			if payment.Status == PaymentStatusPaid {
				addMoney(&paidExpense, payment.Amount)
			}
			if payment.IsOverdue {
				addMoney(&overdueExpense, payment.Amount)
				overdueCount++
			}
		}
	}

	plannedBalance := new(big.Rat).Sub(&plannedIncome, &plannedExpense)
	paidBalance := new(big.Rat).Sub(&paidIncome, &paidExpense)

	return &ProjectFinancialSummary{
		ProjectID:       projectID,
		PlannedIncome:   formatMoney(&plannedIncome),
		PlannedExpense:  formatMoney(&plannedExpense),
		PlannedBalance:  formatMoney(plannedBalance),
		PaidIncome:      formatMoney(&paidIncome),
		PaidExpense:     formatMoney(&paidExpense),
		PaidBalance:     formatMoney(paidBalance),
		OverdueIncome:   formatMoney(&overdueIncome),
		OverdueExpense:  formatMoney(&overdueExpense),
		OverdueCount:    overdueCount,
	}
}

func addMoney(sum *big.Rat, value string) {
	var amount big.Rat
	if _, ok := amount.SetString(strings.TrimSpace(value)); !ok {
		return
	}
	sum.Add(sum, &amount)
}

func formatMoney(value *big.Rat) string {
	if value == nil {
		return "0.00"
	}
	return value.FloatString(2)
}

func IsNotFound(err error) bool {
	return errors.Is(err, ErrPaymentNotFound) || errors.Is(err, ErrProjectNotFound) || errors.Is(err, ErrStageNotFound)
}
