package report

import (
	"net/http"
	"strings"

	"github.com/Hubcher/project-management/gateway/internal/adapters/rest/httpx"
	"github.com/Hubcher/project-management/gateway/internal/adapters/rest/middleware"
	"github.com/Hubcher/project-management/gateway/internal/core"
)

// NewCreateCommentHandler godoc
// @Summary Create daily report comment
// @Description Adds a comment to a daily report. Accessible to users who can read the report, including the author, administrators, and reviewers.
// @Tags report-comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param reportId path string true "Daily report ID"
// @Param request body CreateDailyReportCommentRequest true "Daily report comment payload"
// @Success 201 {object} DailyReportCommentResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/reports/{reportId}/comments [post]
func NewCreateCommentHandler(reports ReportHandlerService, projects core.ProjectDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reportID := strings.TrimSpace(r.PathValue("reportId"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		if _, err := loadReadableReport(r.Context(), reports, projects, authUser, reportID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		var req CreateDailyReportCommentRequest
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		comment, err := reports.CreateComment(r.Context(), toCreateCommentInput(reportID, authUser.UserID, req))
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusCreated, toCommentResponse(comment))
	}
}

// NewListCommentsHandler godoc
// @Summary List daily report comments
// @Description Returns comments attached to a daily report. Accessible to users who can read the report.
// @Tags report-comments
// @Produce json
// @Security BearerAuth
// @Param reportId path string true "Daily report ID"
// @Success 200 {object} ListDailyReportCommentsResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/reports/{reportId}/comments [get]
func NewListCommentsHandler(reports ReportHandlerService, projects core.ProjectDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reportID := strings.TrimSpace(r.PathValue("reportId"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		if _, err := loadReadableReport(r.Context(), reports, projects, authUser, reportID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		comments, err := reports.ListComments(r.Context(), reportID)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, toListCommentsResponse(comments))
	}
}

// NewGetCommentHandler godoc
// @Summary Get daily report comment
// @Description Returns a daily report comment by ID. Accessible to users who can read the parent report.
// @Tags report-comments
// @Produce json
// @Security BearerAuth
// @Param id path string true "Daily report comment ID"
// @Success 200 {object} DailyReportCommentResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/report-comments/{id} [get]
func NewGetCommentHandler(reports ReportHandlerService, projects core.ProjectDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.PathValue("id"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		comment, err := reports.GetComment(r.Context(), id)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if _, err = loadReadableReport(r.Context(), reports, projects, authUser, comment.ReportID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, toCommentResponse(comment))
	}
}

// NewUpdateCommentHandler godoc
// @Summary Update daily report comment
// @Description Updates a daily report comment. Accessible to the comment author and administrators.
// @Tags report-comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Daily report comment ID"
// @Param request body UpdateDailyReportCommentRequest true "Daily report comment update payload"
// @Success 200 {object} DailyReportCommentResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/report-comments/{id} [put]
func NewUpdateCommentHandler(reports ReportHandlerService, projects core.ProjectDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.PathValue("id"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		comment, err := reports.GetComment(r.Context(), id)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if _, err = loadReadableReport(r.Context(), reports, projects, authUser, comment.ReportID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if authUser.Role != core.RoleAdmin && comment.AuthorUserID != authUser.UserID {
			httpx.WriteError(w, http.StatusForbidden, "access denied")
			return
		}

		var req UpdateDailyReportCommentRequest
		if err = httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		updated, err := reports.UpdateComment(r.Context(), toUpdateCommentInput(id, req))
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, toCommentResponse(updated))
	}
}

// NewDeleteCommentHandler godoc
// @Summary Delete daily report comment
// @Description Deletes a daily report comment. Accessible to the comment author and administrators.
// @Tags report-comments
// @Security BearerAuth
// @Param id path string true "Daily report comment ID"
// @Success 204 {string} string "No Content"
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/report-comments/{id} [delete]
func NewDeleteCommentHandler(reports ReportHandlerService, projects core.ProjectDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.PathValue("id"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		comment, err := reports.GetComment(r.Context(), id)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if authUser.Role != core.RoleAdmin && comment.AuthorUserID != authUser.UserID {
			httpx.WriteError(w, http.StatusForbidden, "access denied")
			return
		}
		if _, err = loadReadableReport(r.Context(), reports, projects, authUser, comment.ReportID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if err = reports.DeleteComment(r.Context(), id); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
