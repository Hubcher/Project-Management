package report

import (
	"strings"

	"github.com/Hubcher/project-management/gateway/internal/core"
)

func toCreateReportInput(req CreateDailyReportRequest, userID string) core.CreateDailyReportInput {
	return core.CreateDailyReportInput{UserID: userID, ReportDate: req.ReportDate, Status: req.Status, Summary: req.Summary}
}

func toUpdateReportInput(id string, status, summary string) core.UpdateDailyReportInput {
	return core.UpdateDailyReportInput{ID: id, Status: status, Summary: summary}
}

func toCreateEntryInput(reportID string, req CreateDailyReportEntryRequest) core.CreateDailyReportEntryInput {
	return core.CreateDailyReportEntryInput{ReportID: reportID, ProjectID: req.ProjectID, StageID: req.StageID, WorkType: req.WorkType, Description: req.Description, HoursSpent: req.HoursSpent}
}

func toUpdateEntryInput(id, projectID, stageID, workType, description, hoursSpent string) core.UpdateDailyReportEntryInput {
	return core.UpdateDailyReportEntryInput{ID: id, ProjectID: projectID, StageID: stageID, WorkType: workType, Description: description, HoursSpent: hoursSpent}
}

func toCreateCommentInput(reportID, authorUserID string, req CreateDailyReportCommentRequest) core.CreateDailyReportCommentInput {
	return core.CreateDailyReportCommentInput{ReportID: reportID, AuthorUserID: authorUserID, Comment: req.Comment}
}

func toUpdateCommentInput(id string, req UpdateDailyReportCommentRequest) core.UpdateDailyReportCommentInput {
	return core.UpdateDailyReportCommentInput{ID: id, Comment: req.Comment}
}

func toReportResponse(report *core.DailyReport) DailyReportResponse {
	return DailyReportResponse{ID: report.ID, UserID: report.UserID, ReportDate: report.ReportDate, Status: report.Status, TotalHours: report.TotalHours, Summary: report.Summary, CreatedAt: report.CreatedAt, UpdatedAt: report.UpdatedAt}
}

func toEntryResponse(entry *core.DailyReportEntry) DailyReportEntryResponse {
	return DailyReportEntryResponse{ID: entry.ID, ReportID: entry.ReportID, ProjectID: entry.ProjectID, StageID: entry.StageID, WorkType: entry.WorkType, Description: entry.Description, HoursSpent: entry.HoursSpent, CreatedAt: entry.CreatedAt, UpdatedAt: entry.UpdatedAt}
}

func toCommentResponse(comment *core.DailyReportComment) DailyReportCommentResponse {
	return DailyReportCommentResponse{ID: comment.ID, ReportID: comment.ReportID, AuthorUserID: comment.AuthorUserID, Comment: comment.Comment, CreatedAt: comment.CreatedAt}
}

func toListReportsResponse(reports []core.DailyReport) ListDailyReportsResponse {
	resp := ListDailyReportsResponse{Reports: make([]DailyReportResponse, 0, len(reports))}
	for i := range reports {
		report := reports[i]
		resp.Reports = append(resp.Reports, toReportResponse(&report))
	}
	return resp
}

func toListEntriesResponse(entries []core.DailyReportEntry) ListDailyReportEntriesResponse {
	resp := ListDailyReportEntriesResponse{Entries: make([]DailyReportEntryResponse, 0, len(entries))}
	for i := range entries {
		entry := entries[i]
		resp.Entries = append(resp.Entries, toEntryResponse(&entry))
	}
	return resp
}

func toListCommentsResponse(comments []core.DailyReportComment) ListDailyReportCommentsResponse {
	resp := ListDailyReportCommentsResponse{Comments: make([]DailyReportCommentResponse, 0, len(comments))}
	for i := range comments {
		comment := comments[i]
		resp.Comments = append(resp.Comments, toCommentResponse(&comment))
	}
	return resp
}

func mergeOptionalString(current string, provided *string) string {
	if provided == nil {
		return current
	}
	return strings.TrimSpace(*provided)
}