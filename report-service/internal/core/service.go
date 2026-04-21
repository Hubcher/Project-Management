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
	repo ReportRepository
}

func NewService(repo ReportRepository) ReportService {
	return &service{repo: repo}
}

func (s *service) CreateReport(ctx context.Context, input CreateDailyReportInput) (*DailyReport, error) {
	input = normalizeCreateReportInput(input)
	if err := validateCreateReportInput(input); err != nil {
		return nil, err
	}
	return s.repo.CreateReport(ctx, input)
}

func (s *service) GetReport(ctx context.Context, id string) (*DailyReport, error) {
	id = strings.TrimSpace(id)
	if !isValidUUID(id) {
		return nil, ErrInvalidReport
	}
	return s.repo.GetReport(ctx, id)
}

func (s *service) ListReports(ctx context.Context, filter ListDailyReportsFilter) ([]DailyReport, error) {
	filter = normalizeReportFilter(filter)
	if err := validateReportFilter(filter); err != nil {
		return nil, err
	}
	return s.repo.ListReports(ctx, filter)
}

func (s *service) UpdateReport(ctx context.Context, input UpdateDailyReportInput) (*DailyReport, error) {
	input = normalizeUpdateReportInput(input)
	if !isValidUUID(input.ID) || !isValidReportStatus(input.Status) {
		return nil, ErrInvalidReport
	}
	return s.repo.UpdateReport(ctx, input)
}

func (s *service) DeleteReport(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if !isValidUUID(id) {
		return ErrInvalidReport
	}
	return s.repo.DeleteReport(ctx, id)
}

func (s *service) CreateEntry(ctx context.Context, input CreateDailyReportEntryInput) (*DailyReportEntry, error) {
	input = normalizeCreateEntryInput(input)
	if !isValidUUID(input.ReportID) {
		return nil, ErrInvalidEntry
	}
	if _, err := s.repo.GetReport(ctx, input.ReportID); err != nil {
		return nil, err
	}
	if err := validateEntryInput(input.ProjectID, input.StageID, input.WorkType, input.Description, input.HoursSpent); err != nil {
		return nil, err
	}
	return s.repo.CreateEntry(ctx, input)
}

func (s *service) GetEntry(ctx context.Context, id string) (*DailyReportEntry, error) {
	id = strings.TrimSpace(id)
	if !isValidUUID(id) {
		return nil, ErrInvalidEntry
	}
	return s.repo.GetEntry(ctx, id)
}

func (s *service) ListEntries(ctx context.Context, reportID string) ([]DailyReportEntry, error) {
	reportID = strings.TrimSpace(reportID)
	if !isValidUUID(reportID) {
		return nil, ErrInvalidEntry
	}
	if _, err := s.repo.GetReport(ctx, reportID); err != nil {
		return nil, err
	}
	return s.repo.ListEntries(ctx, reportID)
}

func (s *service) UpdateEntry(ctx context.Context, input UpdateDailyReportEntryInput) (*DailyReportEntry, error) {
	input = normalizeUpdateEntryInput(input)
	if !isValidUUID(input.ID) {
		return nil, ErrInvalidEntry
	}
	if err := validateEntryInput(input.ProjectID, input.StageID, input.WorkType, input.Description, input.HoursSpent); err != nil {
		return nil, err
	}
	return s.repo.UpdateEntry(ctx, input)
}

func (s *service) DeleteEntry(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if !isValidUUID(id) {
		return ErrInvalidEntry
	}
	return s.repo.DeleteEntry(ctx, id)
}

func (s *service) CreateComment(ctx context.Context, input CreateDailyReportCommentInput) (*DailyReportComment, error) {
	input = normalizeCreateCommentInput(input)
	if !isValidUUID(input.ReportID) {
		return nil, ErrInvalidComment
	}
	if _, err := s.repo.GetReport(ctx, input.ReportID); err != nil {
		return nil, err
	}
	if err := validateCommentInput(input.AuthorUserID, input.Comment); err != nil {
		return nil, err
	}
	return s.repo.CreateComment(ctx, input)
}

func (s *service) GetComment(ctx context.Context, id string) (*DailyReportComment, error) {
	id = strings.TrimSpace(id)
	if !isValidUUID(id) {
		return nil, ErrInvalidComment
	}
	return s.repo.GetComment(ctx, id)
}

func (s *service) ListComments(ctx context.Context, reportID string) ([]DailyReportComment, error) {
	reportID = strings.TrimSpace(reportID)
	if !isValidUUID(reportID) {
		return nil, ErrInvalidComment
	}
	if _, err := s.repo.GetReport(ctx, reportID); err != nil {
		return nil, err
	}
	return s.repo.ListComments(ctx, reportID)
}

func (s *service) UpdateComment(ctx context.Context, input UpdateDailyReportCommentInput) (*DailyReportComment, error) {
	input = normalizeUpdateCommentInput(input)
	if !isValidUUID(input.ID) || strings.TrimSpace(input.Comment) == "" {
		return nil, ErrInvalidComment
	}
	return s.repo.UpdateComment(ctx, input)
}

func (s *service) DeleteComment(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if !isValidUUID(id) {
		return ErrInvalidComment
	}
	return s.repo.DeleteComment(ctx, id)
}

func validateCreateReportInput(input CreateDailyReportInput) error {
	if !isValidUUID(input.UserID) || input.ReportDate == nil || !isValidReportStatus(input.Status) {
		return ErrInvalidReport
	}
	return nil
}

func validateReportFilter(filter ListDailyReportsFilter) error {
	if filter.UserID != "" && !isValidUUID(filter.UserID) {
		return ErrInvalidReport
	}
	if filter.Status != "" && !isValidReportStatus(filter.Status) {
		return ErrInvalidReport
	}
	if !isValidDateRange(filter.DateFrom, filter.DateTo) {
		return ErrInvalidReport
	}
	return nil
}

func validateEntryInput(projectID, stageID, workType, description, hoursSpent string) error {
	if !isValidUUID(projectID) || strings.TrimSpace(workType) == "" || strings.TrimSpace(description) == "" {
		return ErrInvalidEntry
	}
	if strings.TrimSpace(stageID) != "" && !isValidUUID(stageID) {
		return ErrInvalidEntry
	}
	if !isPositiveMoney(hoursSpent) {
		return ErrInvalidEntry
	}
	return nil
}

func validateCommentInput(authorUserID, comment string) error {
	if !isValidUUID(authorUserID) || strings.TrimSpace(comment) == "" {
		return ErrInvalidComment
	}
	return nil
}

func normalizeCreateReportInput(input CreateDailyReportInput) CreateDailyReportInput {
	input.UserID = strings.TrimSpace(input.UserID)
	input.ReportDate = normalizeDate(input.ReportDate)
	input.Status = normalizeReportStatus(input.Status, true)
	input.Summary = strings.TrimSpace(input.Summary)
	return input
}

func normalizeUpdateReportInput(input UpdateDailyReportInput) UpdateDailyReportInput {
	input.ID = strings.TrimSpace(input.ID)
	input.Status = normalizeReportStatus(input.Status, false)
	input.Summary = strings.TrimSpace(input.Summary)
	return input
}

func normalizeReportFilter(filter ListDailyReportsFilter) ListDailyReportsFilter {
	filter.UserID = strings.TrimSpace(filter.UserID)
	filter.Status = normalizeReportStatus(filter.Status, false)
	filter.DateFrom = normalizeDate(filter.DateFrom)
	filter.DateTo = normalizeDate(filter.DateTo)
	return filter
}

func normalizeCreateEntryInput(input CreateDailyReportEntryInput) CreateDailyReportEntryInput {
	input.ReportID = strings.TrimSpace(input.ReportID)
	input.ProjectID = strings.TrimSpace(input.ProjectID)
	input.StageID = strings.TrimSpace(input.StageID)
	input.WorkType = strings.TrimSpace(input.WorkType)
	input.Description = strings.TrimSpace(input.Description)
	input.HoursSpent = normalizeMoney(input.HoursSpent)
	return input
}

func normalizeUpdateEntryInput(input UpdateDailyReportEntryInput) UpdateDailyReportEntryInput {
	input.ID = strings.TrimSpace(input.ID)
	input.ProjectID = strings.TrimSpace(input.ProjectID)
	input.StageID = strings.TrimSpace(input.StageID)
	input.WorkType = strings.TrimSpace(input.WorkType)
	input.Description = strings.TrimSpace(input.Description)
	input.HoursSpent = normalizeMoney(input.HoursSpent)
	return input
}

func normalizeCreateCommentInput(input CreateDailyReportCommentInput) CreateDailyReportCommentInput {
	input.ReportID = strings.TrimSpace(input.ReportID)
	input.AuthorUserID = strings.TrimSpace(input.AuthorUserID)
	input.Comment = strings.TrimSpace(input.Comment)
	return input
}

func normalizeUpdateCommentInput(input UpdateDailyReportCommentInput) UpdateDailyReportCommentInput {
	input.ID = strings.TrimSpace(input.ID)
	input.Comment = strings.TrimSpace(input.Comment)
	return input
}

func normalizeReportStatus(status ReportStatus, allowDefault bool) ReportStatus {
	status = ReportStatus(strings.ToLower(strings.TrimSpace(string(status))))
	if status == "" && allowDefault {
		return ReportStatusDraft
	}
	return status
}

func normalizeDate(value *time.Time) *time.Time {
	if value == nil || value.IsZero() {
		return nil
	}
	day := value.UTC()
	normalized := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	return &normalized
}

func normalizeMoney(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "0"
	}
	return value
}

func isValidUUID(value string) bool {
	_, err := uuid.Parse(strings.TrimSpace(value))
	return err == nil
}

func isPositiveMoney(value string) bool {
	var amount big.Rat
	if _, ok := amount.SetString(strings.TrimSpace(value)); !ok {
		return false
	}
	return amount.Sign() > 0
}

func isValidDateRange(start, end *time.Time) bool {
	if start == nil || end == nil {
		return true
	}
	return !end.Before(*start)
}

func isValidReportStatus(status ReportStatus) bool {
	switch status {
	case ReportStatusDraft, ReportStatusSubmitted, ReportStatusApproved, ReportStatusRejected:
		return true
	default:
		return false
	}
}

func IsNotFound(err error) bool {
	return errors.Is(err, ErrReportNotFound) || errors.Is(err, ErrEntryNotFound) || errors.Is(err, ErrCommentNotFound)
}