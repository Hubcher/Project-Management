package core

import "time"

type ReportStatus string

const (
	ReportStatusDraft     ReportStatus = "draft"
	ReportStatusSubmitted ReportStatus = "submitted"
	ReportStatusApproved  ReportStatus = "approved"
	ReportStatusRejected  ReportStatus = "rejected"
)

type DailyReport struct {
	ID         string
	UserID     string
	ReportDate time.Time
	Status     ReportStatus
	TotalHours string
	Summary    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ListDailyReportsFilter struct {
	UserID   string
	Status   ReportStatus
	DateFrom *time.Time
	DateTo   *time.Time
}

type CreateDailyReportInput struct {
	UserID     string
	ReportDate *time.Time
	Status     ReportStatus
	Summary    string
}

type UpdateDailyReportInput struct {
	ID      string
	Status  ReportStatus
	Summary string
}

type DailyReportEntry struct {
	ID          string
	ReportID    string
	ProjectID   string
	StageID     string
	WorkType    string
	Description string
	HoursSpent  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateDailyReportEntryInput struct {
	ReportID    string
	ProjectID   string
	StageID     string
	WorkType    string
	Description string
	HoursSpent  string
}

type UpdateDailyReportEntryInput struct {
	ID          string
	ProjectID   string
	StageID     string
	WorkType    string
	Description string
	HoursSpent  string
}

type DailyReportComment struct {
	ID           string
	ReportID     string
	AuthorUserID string
	Comment      string
	CreatedAt    time.Time
}

type CreateDailyReportCommentInput struct {
	ReportID     string
	AuthorUserID string
	Comment      string
}

type UpdateDailyReportCommentInput struct {
	ID      string
	Comment string
}