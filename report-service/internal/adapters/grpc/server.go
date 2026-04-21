package grpc

import (
	"context"
	"errors"
	"time"

	reportpb "github.com/Hubcher/project-management/contracts/gen/go/report"
	"github.com/Hubcher/project-management/report-service/internal/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	reportpb.UnimplementedReportServiceServer
	service core.ReportService
}

func NewServer(service core.ReportService) *Server {
	return &Server{service: service}
}

func (s *Server) Ping(_ context.Context, _ *reportpb.Empty) (*reportpb.Empty, error) {
	return &reportpb.Empty{}, nil
}

func (s *Server) CreateReport(ctx context.Context, req *reportpb.CreateDailyReportRequest) (*reportpb.DailyReport, error) {
	reportDate, err := parseDate(req.ReportDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid report_date")
	}
	report, err := s.service.CreateReport(ctx, core.CreateDailyReportInput{
		UserID:     req.UserId,
		ReportDate: reportDate,
		Status:     core.ReportStatus(req.Status),
		Summary:    req.Summary,
	})
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoReport(report), nil
}

func (s *Server) GetReport(ctx context.Context, req *reportpb.GetDailyReportRequest) (*reportpb.DailyReport, error) {
	report, err := s.service.GetReport(ctx, req.Id)
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoReport(report), nil
}

func (s *Server) ListReports(ctx context.Context, req *reportpb.ListDailyReportsRequest) (*reportpb.ListDailyReportsResponse, error) {
	dateFrom, err := parseDate(req.DateFrom)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid date_from")
	}
	dateTo, err := parseDate(req.DateTo)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid date_to")
	}
	reports, err := s.service.ListReports(ctx, core.ListDailyReportsFilter{
		UserID:   req.UserId,
		Status:   core.ReportStatus(req.Status),
		DateFrom: dateFrom,
		DateTo:   dateTo,
	})
	if err != nil {
		return nil, mapCoreError(err)
	}
	resp := &reportpb.ListDailyReportsResponse{Reports: make([]*reportpb.DailyReport, 0, len(reports))}
	for i := range reports {
		resp.Reports = append(resp.Reports, toProtoReport(&reports[i]))
	}
	return resp, nil
}

func (s *Server) UpdateReport(ctx context.Context, req *reportpb.UpdateDailyReportRequest) (*reportpb.DailyReport, error) {
	report, err := s.service.UpdateReport(ctx, core.UpdateDailyReportInput{
		ID:      req.Id,
		Status:  core.ReportStatus(req.Status),
		Summary: req.Summary,
	})
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoReport(report), nil
}

func (s *Server) DeleteReport(ctx context.Context, req *reportpb.DeleteDailyReportRequest) (*reportpb.Empty, error) {
	if err := s.service.DeleteReport(ctx, req.Id); err != nil {
		return nil, mapCoreError(err)
	}
	return &reportpb.Empty{}, nil
}

func (s *Server) CreateEntry(ctx context.Context, req *reportpb.CreateDailyReportEntryRequest) (*reportpb.DailyReportEntry, error) {
	entry, err := s.service.CreateEntry(ctx, core.CreateDailyReportEntryInput{
		ReportID:    req.ReportId,
		ProjectID:   req.ProjectId,
		StageID:     req.StageId,
		WorkType:    req.WorkType,
		Description: req.Description,
		HoursSpent:  req.HoursSpent,
	})
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoEntry(entry), nil
}

func (s *Server) GetEntry(ctx context.Context, req *reportpb.GetDailyReportEntryRequest) (*reportpb.DailyReportEntry, error) {
	entry, err := s.service.GetEntry(ctx, req.Id)
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoEntry(entry), nil
}

func (s *Server) ListEntries(ctx context.Context, req *reportpb.ListDailyReportEntriesRequest) (*reportpb.ListDailyReportEntriesResponse, error) {
	entries, err := s.service.ListEntries(ctx, req.ReportId)
	if err != nil {
		return nil, mapCoreError(err)
	}
	resp := &reportpb.ListDailyReportEntriesResponse{Entries: make([]*reportpb.DailyReportEntry, 0, len(entries))}
	for i := range entries {
		resp.Entries = append(resp.Entries, toProtoEntry(&entries[i]))
	}
	return resp, nil
}

func (s *Server) UpdateEntry(ctx context.Context, req *reportpb.UpdateDailyReportEntryRequest) (*reportpb.DailyReportEntry, error) {
	entry, err := s.service.UpdateEntry(ctx, core.UpdateDailyReportEntryInput{
		ID:          req.Id,
		ProjectID:   req.ProjectId,
		StageID:     req.StageId,
		WorkType:    req.WorkType,
		Description: req.Description,
		HoursSpent:  req.HoursSpent,
	})
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoEntry(entry), nil
}

func (s *Server) DeleteEntry(ctx context.Context, req *reportpb.DeleteDailyReportEntryRequest) (*reportpb.Empty, error) {
	if err := s.service.DeleteEntry(ctx, req.Id); err != nil {
		return nil, mapCoreError(err)
	}
	return &reportpb.Empty{}, nil
}

func (s *Server) CreateComment(ctx context.Context, req *reportpb.CreateDailyReportCommentRequest) (*reportpb.DailyReportComment, error) {
	comment, err := s.service.CreateComment(ctx, core.CreateDailyReportCommentInput{
		ReportID:     req.ReportId,
		AuthorUserID: req.AuthorUserId,
		Comment:      req.Comment,
	})
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoComment(comment), nil
}

func (s *Server) GetComment(ctx context.Context, req *reportpb.GetDailyReportCommentRequest) (*reportpb.DailyReportComment, error) {
	comment, err := s.service.GetComment(ctx, req.Id)
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoComment(comment), nil
}

func (s *Server) ListComments(ctx context.Context, req *reportpb.ListDailyReportCommentsRequest) (*reportpb.ListDailyReportCommentsResponse, error) {
	comments, err := s.service.ListComments(ctx, req.ReportId)
	if err != nil {
		return nil, mapCoreError(err)
	}
	resp := &reportpb.ListDailyReportCommentsResponse{Comments: make([]*reportpb.DailyReportComment, 0, len(comments))}
	for i := range comments {
		resp.Comments = append(resp.Comments, toProtoComment(&comments[i]))
	}
	return resp, nil
}

func (s *Server) UpdateComment(ctx context.Context, req *reportpb.UpdateDailyReportCommentRequest) (*reportpb.DailyReportComment, error) {
	comment, err := s.service.UpdateComment(ctx, core.UpdateDailyReportCommentInput{
		ID:      req.Id,
		Comment: req.Comment,
	})
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoComment(comment), nil
}

func (s *Server) DeleteComment(ctx context.Context, req *reportpb.DeleteDailyReportCommentRequest) (*reportpb.Empty, error) {
	if err := s.service.DeleteComment(ctx, req.Id); err != nil {
		return nil, mapCoreError(err)
	}
	return &reportpb.Empty{}, nil
}

func parseDate(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, err
	}
	day := parsed.UTC()
	normalized := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	return &normalized, nil
}

func toProtoReport(report *core.DailyReport) *reportpb.DailyReport {
	return &reportpb.DailyReport{
		ID:         report.ID,
		UserId:     report.UserID,
		ReportDate: formatDate(&report.ReportDate),
		Status:     string(report.Status),
		TotalHours: report.TotalHours,
		Summary:    report.Summary,
		CreatedAt:  formatTimestamp(&report.CreatedAt),
		UpdatedAt:  formatTimestamp(&report.UpdatedAt),
	}
}

func toProtoEntry(entry *core.DailyReportEntry) *reportpb.DailyReportEntry {
	return &reportpb.DailyReportEntry{
		ID:          entry.ID,
		ReportId:    entry.ReportID,
		ProjectId:   entry.ProjectID,
		StageId:     entry.StageID,
		WorkType:    entry.WorkType,
		Description: entry.Description,
		HoursSpent:  entry.HoursSpent,
		CreatedAt:   formatTimestamp(&entry.CreatedAt),
		UpdatedAt:   formatTimestamp(&entry.UpdatedAt),
	}
}

func toProtoComment(comment *core.DailyReportComment) *reportpb.DailyReportComment {
	return &reportpb.DailyReportComment{
		ID:           comment.ID,
		ReportId:     comment.ReportID,
		AuthorUserId: comment.AuthorUserID,
		Comment:      comment.Comment,
		CreatedAt:    formatTimestamp(&comment.CreatedAt),
	}
}

func formatDate(value *time.Time) string {
	if value == nil || value.IsZero() {
		return ""
	}
	return value.UTC().Format("2006-01-02")
}

func formatTimestamp(value *time.Time) string {
	if value == nil || value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339)
}

func mapCoreError(err error) error {
	switch {
	case errors.Is(err, core.ErrInvalidReport), errors.Is(err, core.ErrInvalidEntry), errors.Is(err, core.ErrInvalidComment):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, core.ErrReportNotFound), errors.Is(err, core.ErrEntryNotFound), errors.Is(err, core.ErrCommentNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, core.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}