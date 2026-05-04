package report

import (
	"context"
	"net/http"
	"strings"

	"github.com/Hubcher/project-management/gateway/internal/adapters/rest/httpx"
	"github.com/Hubcher/project-management/gateway/internal/adapters/rest/middleware"
	"github.com/Hubcher/project-management/gateway/internal/core"
)

const reviewScope = "review"

// NewCreateReportHandler godoc
// @Summary Create daily report
// @Description Creates a daily report. Regular users can create reports only for themselves, administrators may specify another user via `user_id`.
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateDailyReportRequest true "Daily report payload"
// @Success 201 {object} DailyReportResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/reports [post]
func NewCreateReportHandler(reports ReportHandlerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		var req CreateDailyReportRequest
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		ownerID, err := resolveReportOwner(authUser, req.UserID)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if authUser.Role != core.RoleAdmin && isReviewStatus(req.Status) {
			httpx.WriteError(w, http.StatusForbidden, "access denied")
			return
		}

		report, err := reports.CreateReport(r.Context(), toCreateReportInput(req, ownerID))
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusCreated, toReportResponse(report))
	}
}

// NewGetReportHandler godoc
// @Summary Get daily report
// @Description Returns a daily report by ID. Accessible to the report author, administrators, and project managers reviewing related project work.
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Param id path string true "Daily report ID"
// @Success 200 {object} DailyReportResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/reports/{id} [get]
func NewGetReportHandler(reports ReportHandlerService, projects core.ProjectDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reportID := strings.TrimSpace(r.PathValue("id"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		report, err := loadReadableReport(r.Context(), reports, projects, authUser, reportID)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, toReportResponse(report))
	}
}

// NewListReportsHandler godoc
// @Summary List daily reports
// @Description Returns daily reports visible to the caller. Regular users list their own reports, managers may request `scope=review` to see reviewable reports, and administrators can list all reports.
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Param scope query string false "Listing scope. Use `self` for own reports or `review` for manager review queue."
// @Param user_id query string false "User ID filter. Available to administrators, or must match the current user for `scope=self`."
// @Param status query string false "Report status filter"
// @Param date_from query string false "Lower bound for report date in YYYY-MM-DD format"
// @Param date_to query string false "Upper bound for report date in YYYY-MM-DD format"
// @Success 200 {object} ListDailyReportsResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/reports [get]
func NewListReportsHandler(reports ReportHandlerService, projects core.ProjectDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		filter := core.ListDailyReportsInput{
			UserID:   strings.TrimSpace(r.URL.Query().Get("user_id")),
			Status:   strings.TrimSpace(r.URL.Query().Get("status")),
			DateFrom: strings.TrimSpace(r.URL.Query().Get("date_from")),
			DateTo:   strings.TrimSpace(r.URL.Query().Get("date_to")),
		}
		scope := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("scope")))
		if scope == "" {
			scope = "self"
		}

		var (
			list []core.DailyReport
			err  error
		)

		switch {
		case authUser.Role == core.RoleAdmin:
			list, err = reports.ListReports(r.Context(), filter)
		case scope == reviewScope && authUser.Role == core.RoleManager:
			list, err = listReviewableReports(r.Context(), reports, projects, authUser, filter)
		case scope == reviewScope:
			httpx.WriteError(w, http.StatusForbidden, "access denied")
			return
		case scope == "self":
			if filter.UserID != "" && filter.UserID != authUser.UserID {
				httpx.WriteError(w, http.StatusForbidden, "access denied")
				return
			}
			filter.UserID = authUser.UserID
			list, err = reports.ListReports(r.Context(), filter)
		default:
			httpx.WriteError(w, http.StatusBadRequest, "invalid scope")
			return
		}
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, toListReportsResponse(list))
	}
}

// NewUpdateReportHandler godoc
// @Summary Update daily report
// @Description Updates a daily report. Authors may edit draft or submitted report fields they own, while administrators and reviewers may set review statuses such as `approved` or `rejected`.
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Daily report ID"
// @Param request body UpdateDailyReportRequest true "Daily report update payload"
// @Success 200 {object} DailyReportResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/reports/{id} [put]
func NewUpdateReportHandler(reports ReportHandlerService, projects core.ProjectDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reportID := strings.TrimSpace(r.PathValue("id"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		report, err := reports.GetReport(r.Context(), reportID)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if err = ensureReadableReport(r.Context(), reports, projects, authUser, report); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		var req UpdateDailyReportRequest
		if err = httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		if req.Status == nil && req.Summary == nil {
			httpx.WriteJSON(w, http.StatusOK, toReportResponse(report))
			return
		}

		isAdmin := authUser.Role == core.RoleAdmin
		isOwner := report.UserID == authUser.UserID
		canReview, err := isReportReviewer(r.Context(), reports, projects, authUser, report)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		statusValue := report.Status
		summaryValue := report.Summary
		if req.Summary != nil {
			if !isAdmin && !isOwner {
				httpx.WriteError(w, http.StatusForbidden, "access denied")
				return
			}
			summaryValue = strings.TrimSpace(*req.Summary)
		}
		if req.Status != nil {
			nextStatus := strings.TrimSpace(*req.Status)
			if isReviewStatus(nextStatus) {
				if !isAdmin && !canReview {
					httpx.WriteError(w, http.StatusForbidden, "access denied")
					return
				}
			} else if !isAdmin && !isOwner {
				httpx.WriteError(w, http.StatusForbidden, "access denied")
				return
			}
			statusValue = nextStatus
		}

		updated, err := reports.UpdateReport(r.Context(), toUpdateReportInput(reportID, statusValue, summaryValue))
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, toReportResponse(updated))
	}
}

// NewDeleteReportHandler godoc
// @Summary Delete daily report
// @Description Deletes a daily report. Accessible to the report author and administrators.
// @Tags reports
// @Security BearerAuth
// @Param id path string true "Daily report ID"
// @Success 204 {string} string "No Content"
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/reports/{id} [delete]
func NewDeleteReportHandler(reports ReportHandlerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reportID := strings.TrimSpace(r.PathValue("id"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		if _, err := loadOwnedReport(r.Context(), reports, authUser, reportID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if err := reports.DeleteReport(r.Context(), reportID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func loadReadableReport(ctx context.Context, reports ReportHandlerService, projects core.ProjectDirectory, authUser core.AuthUser, reportID string) (*core.DailyReport, error) {
	report, err := reports.GetReport(ctx, reportID)
	if err != nil {
		return nil, err
	}
	if err = ensureReadableReport(ctx, reports, projects, authUser, report); err != nil {
		return nil, err
	}
	return report, nil
}

func ensureReadableReport(ctx context.Context, reports ReportHandlerService, projects core.ProjectDirectory, authUser core.AuthUser, report *core.DailyReport) error {
	if authUser.Role == core.RoleAdmin || report.UserID == authUser.UserID {
		return nil
	}
	allowed, err := isReportReviewer(ctx, reports, projects, authUser, report)
	if err != nil {
		return err
	}
	if !allowed {
		return core.NewStatusError(http.StatusForbidden, "access denied")
	}
	return nil
}

func loadOwnedReport(ctx context.Context, reports ReportHandlerService, authUser core.AuthUser, reportID string) (*core.DailyReport, error) {
	report, err := reports.GetReport(ctx, reportID)
	if err != nil {
		return nil, err
	}
	if authUser.Role == core.RoleAdmin || report.UserID == authUser.UserID {
		return report, nil
	}
	return nil, core.NewStatusError(http.StatusForbidden, "access denied")
}

func loadEditableReport(ctx context.Context, reports ReportHandlerService, authUser core.AuthUser, reportID string) (*core.DailyReport, error) {
	report, err := loadOwnedReport(ctx, reports, authUser, reportID)
	if err != nil {
		return nil, err
	}
	if authUser.Role != core.RoleAdmin && (report.Status == "submitted" || report.Status == "approved") {
		return nil, core.NewStatusError(http.StatusForbidden, "report is locked for editing")
	}
	return report, nil
}

func resolveReportOwner(authUser core.AuthUser, requestedUserID string) (string, error) {
	requestedUserID = strings.TrimSpace(requestedUserID)
	if authUser.Role == core.RoleAdmin {
		if requestedUserID == "" {
			return authUser.UserID, nil
		}
		return requestedUserID, nil
	}
	if requestedUserID != "" && requestedUserID != authUser.UserID {
		return "", core.NewStatusError(http.StatusForbidden, "access denied")
	}
	return authUser.UserID, nil
}

func listReviewableReports(ctx context.Context, reports ReportHandlerService, projects core.ProjectDirectory, authUser core.AuthUser, filter core.ListDailyReportsInput) ([]core.DailyReport, error) {
	candidates, err := reports.ListReports(ctx, filter)
	if err != nil {
		return nil, err
	}
	result := make([]core.DailyReport, 0, len(candidates))
	for i := range candidates {
		report := candidates[i]
		allowed, reviewErr := isReportReviewer(ctx, reports, projects, authUser, &report)
		if reviewErr != nil {
			return nil, reviewErr
		}
		if allowed {
			result = append(result, report)
		}
	}
	return result, nil
}

func isReportReviewer(ctx context.Context, reports ReportHandlerService, projects core.ProjectDirectory, authUser core.AuthUser, report *core.DailyReport) (bool, error) {
	if authUser.Role == core.RoleAdmin {
		return true, nil
	}
	if report.UserID == authUser.UserID {
		return false, nil
	}
	entries, err := reports.ListEntries(ctx, report.ID)
	if err != nil {
		return false, err
	}
	seenProjects := make(map[string]struct{}, len(entries))
	for i := range entries {
		projectID := entries[i].ProjectID
		if _, exists := seenProjects[projectID]; exists {
			continue
		}
		seenProjects[projectID] = struct{}{}
		project, getErr := projects.GetProject(ctx, projectID)
		if getErr != nil {
			return false, getErr
		}
		if authUser.Role == core.RoleManager && project.ManagerID == authUser.UserID {
			return true, nil
		}
	}
	return false, nil
}

func ensureReportOwnerAssignedToProject(ctx context.Context, projects core.ProjectDirectory, report *core.DailyReport, projectID, stageID string) error {
	project, err := projects.GetProject(ctx, strings.TrimSpace(projectID))
	if err != nil {
		return err
	}
	if trimmedStageID := strings.TrimSpace(stageID); trimmedStageID != "" {
		stage, stageErr := projects.GetStage(ctx, trimmedStageID)
		if stageErr != nil {
			return stageErr
		}
		if stage.ProjectID != project.ID {
			return core.NewStatusError(http.StatusBadRequest, "project stage does not belong to project")
		}
	}
	if project.ManagerID == report.UserID {
		return nil
	}
	members, err := projects.ListMembers(ctx, project.ID)
	if err != nil {
		return err
	}
	for _, member := range members {
		if member.UserID == report.UserID && member.IsActive {
			return nil
		}
	}
	return core.NewStatusError(http.StatusBadRequest, "report user is not assigned to project")
}

func isReviewStatus(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "approved", "rejected":
		return true
	default:
		return false
	}
}
