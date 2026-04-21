package core

import "context"

type ReportRepository interface {
	CreateReport(ctx context.Context, input CreateDailyReportInput) (*DailyReport, error)
	GetReport(ctx context.Context, id string) (*DailyReport, error)
	ListReports(ctx context.Context, filter ListDailyReportsFilter) ([]DailyReport, error)
	UpdateReport(ctx context.Context, input UpdateDailyReportInput) (*DailyReport, error)
	DeleteReport(ctx context.Context, id string) error

	CreateEntry(ctx context.Context, input CreateDailyReportEntryInput) (*DailyReportEntry, error)
	GetEntry(ctx context.Context, id string) (*DailyReportEntry, error)
	ListEntries(ctx context.Context, reportID string) ([]DailyReportEntry, error)
	UpdateEntry(ctx context.Context, input UpdateDailyReportEntryInput) (*DailyReportEntry, error)
	DeleteEntry(ctx context.Context, id string) error

	CreateComment(ctx context.Context, input CreateDailyReportCommentInput) (*DailyReportComment, error)
	GetComment(ctx context.Context, id string) (*DailyReportComment, error)
	ListComments(ctx context.Context, reportID string) ([]DailyReportComment, error)
	UpdateComment(ctx context.Context, input UpdateDailyReportCommentInput) (*DailyReportComment, error)
	DeleteComment(ctx context.Context, id string) error
}

type ReportService interface {
	CreateReport(ctx context.Context, input CreateDailyReportInput) (*DailyReport, error)
	GetReport(ctx context.Context, id string) (*DailyReport, error)
	ListReports(ctx context.Context, filter ListDailyReportsFilter) ([]DailyReport, error)
	UpdateReport(ctx context.Context, input UpdateDailyReportInput) (*DailyReport, error)
	DeleteReport(ctx context.Context, id string) error

	CreateEntry(ctx context.Context, input CreateDailyReportEntryInput) (*DailyReportEntry, error)
	GetEntry(ctx context.Context, id string) (*DailyReportEntry, error)
	ListEntries(ctx context.Context, reportID string) ([]DailyReportEntry, error)
	UpdateEntry(ctx context.Context, input UpdateDailyReportEntryInput) (*DailyReportEntry, error)
	DeleteEntry(ctx context.Context, id string) error

	CreateComment(ctx context.Context, input CreateDailyReportCommentInput) (*DailyReportComment, error)
	GetComment(ctx context.Context, id string) (*DailyReportComment, error)
	ListComments(ctx context.Context, reportID string) ([]DailyReportComment, error)
	UpdateComment(ctx context.Context, input UpdateDailyReportCommentInput) (*DailyReportComment, error)
	DeleteComment(ctx context.Context, id string) error
}