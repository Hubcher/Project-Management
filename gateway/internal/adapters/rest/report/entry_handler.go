package report

import (
	"net/http"
	"strings"

	"github.com/Hubcher/project-management/gateway/internal/adapters/rest/httpx"
	"github.com/Hubcher/project-management/gateway/internal/adapters/rest/middleware"
	"github.com/Hubcher/project-management/gateway/internal/core"
)

// NewCreateEntryHandler godoc
// @Summary Create daily report entry
// @Description Adds a worklog entry to a daily report. Accessible to the report author while the report remains editable.
// @Tags report-entries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param reportId path string true "Daily report ID"
// @Param request body CreateDailyReportEntryRequest true "Daily report entry payload"
// @Success 201 {object} DailyReportEntryResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/reports/{reportId}/entries [post]
func NewCreateEntryHandler(reports ReportHandlerService, projects core.ProjectDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reportID := strings.TrimSpace(r.PathValue("reportId"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		report, err := loadEditableReport(r.Context(), reports, authUser, reportID)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		var req CreateDailyReportEntryRequest
		if err = httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err = ensureReportOwnerAssignedToProject(r.Context(), projects, report, req.ProjectID, req.StageID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		entry, err := reports.CreateEntry(r.Context(), toCreateEntryInput(reportID, req))
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusCreated, toEntryResponse(entry))
	}
}

// NewListEntriesHandler godoc
// @Summary List daily report entries
// @Description Returns all entries that belong to a daily report. Accessible to the report author, administrators, and reviewers with access to the report.
// @Tags report-entries
// @Produce json
// @Security BearerAuth
// @Param reportId path string true "Daily report ID"
// @Success 200 {object} ListDailyReportEntriesResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/reports/{reportId}/entries [get]
func NewListEntriesHandler(reports ReportHandlerService, projects core.ProjectDirectory) http.HandlerFunc {
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
		entries, err := reports.ListEntries(r.Context(), reportID)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, toListEntriesResponse(entries))
	}
}

// NewGetEntryHandler godoc
// @Summary Get daily report entry
// @Description Returns a daily report entry by ID. Accessible to the report author, administrators, and reviewers with access to the parent report.
// @Tags report-entries
// @Produce json
// @Security BearerAuth
// @Param id path string true "Daily report entry ID"
// @Success 200 {object} DailyReportEntryResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/report-entries/{id} [get]
func NewGetEntryHandler(reports ReportHandlerService, projects core.ProjectDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.PathValue("id"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		entry, err := reports.GetEntry(r.Context(), id)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if _, err = loadReadableReport(r.Context(), reports, projects, authUser, entry.ReportID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, toEntryResponse(entry))
	}
}

// NewUpdateEntryHandler godoc
// @Summary Update daily report entry
// @Description Updates a daily report entry. Accessible to the report author while the parent report remains editable.
// @Tags report-entries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Daily report entry ID"
// @Param request body UpdateDailyReportEntryRequest true "Daily report entry update payload"
// @Success 200 {object} DailyReportEntryResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/report-entries/{id} [put]
func NewUpdateEntryHandler(reports ReportHandlerService, projects core.ProjectDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.PathValue("id"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		entry, err := reports.GetEntry(r.Context(), id)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		report, err := loadEditableReport(r.Context(), reports, authUser, entry.ReportID)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		var req UpdateDailyReportEntryRequest
		if err = httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		projectID := mergeOptionalString(entry.ProjectID, req.ProjectID)
		stageID := mergeOptionalString(entry.StageID, req.StageID)
		workType := mergeOptionalString(entry.WorkType, req.WorkType)
		description := mergeOptionalString(entry.Description, req.Description)
		hoursSpent := mergeOptionalString(entry.HoursSpent, req.HoursSpent)
		if err = ensureReportOwnerAssignedToProject(r.Context(), projects, report, projectID, stageID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		updated, err := reports.UpdateEntry(r.Context(), toUpdateEntryInput(id, projectID, stageID, workType, description, hoursSpent))
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, toEntryResponse(updated))
	}
}

// NewDeleteEntryHandler godoc
// @Summary Delete daily report entry
// @Description Deletes a daily report entry. Accessible to the report author while the parent report remains editable.
// @Tags report-entries
// @Security BearerAuth
// @Param id path string true "Daily report entry ID"
// @Success 204 {string} string "No Content"
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/report-entries/{id} [delete]
func NewDeleteEntryHandler(reports ReportHandlerService, projects core.ProjectDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.PathValue("id"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		entry, err := reports.GetEntry(r.Context(), id)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if _, err = loadEditableReport(r.Context(), reports, authUser, entry.ReportID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if err = reports.DeleteEntry(r.Context(), id); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
