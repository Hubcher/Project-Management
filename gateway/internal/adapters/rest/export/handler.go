package export

import (
	"mime"
	"net/http"
	"strconv"
	"strings"

	"github.com/Hubcher/project-management/gateway/internal/adapters/rest/httpx"
	"github.com/Hubcher/project-management/gateway/internal/adapters/rest/middleware"
	"github.com/Hubcher/project-management/gateway/internal/core"
)

// NewBuildExportHandler godoc
// @Summary Build management export
// @Description Builds CSV, XLSX or PDF management reports from project and payment calendar data. Regular users receive only their projects; administrators may export the whole portfolio.
// @Tags exports
// @Produce application/octet-stream
// @Security BearerAuth
// @Param reportType path string true "Report type: summary, projects, payments, cashflow, profit or risks"
// @Param format query string false "File format: csv, xlsx or pdf"
// @Param project_id query string false "Project ID filter"
// @Param date_from query string false "Lower payment planned date bound in YYYY-MM-DD format"
// @Param date_to query string false "Upper payment planned date bound in YYYY-MM-DD format"
// @Param group_by query string false "Grouping: day, week or month"
// @Param payment_type query string false "Payment type: income or expense"
// @Param payment_status query string false "Payment status: planned, paid or cancelled"
// @Param overdue_only query bool false "Return only computed overdue payments"
// @Success 200 {file} binary
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/exports/{reportType} [get]
func NewBuildExportHandler(exports ExportHandlerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		input, err := toBuildExportInput(r, authUser)
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		file, err := exports.BuildExport(r.Context(), input)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		w.Header().Set("Content-Type", file.ContentType)
		w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": file.FileName}))
		w.Header().Set("Content-Length", strconv.Itoa(len(file.Data)))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(file.Data)
	}
}

func toBuildExportInput(r *http.Request, authUser core.AuthUser) (core.BuildExportInput, error) {
	values := r.URL.Query()
	overdueOnly := false
	if raw := strings.TrimSpace(values.Get("overdue_only")); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			return core.BuildExportInput{}, err
		}
		overdueOnly = parsed
	}

	input := core.BuildExportInput{
		ReportType:         strings.TrimSpace(r.PathValue("reportType")),
		Format:             strings.TrimSpace(values.Get("format")),
		ProjectID:          strings.TrimSpace(values.Get("project_id")),
		DateFrom:           strings.TrimSpace(values.Get("date_from")),
		DateTo:             strings.TrimSpace(values.Get("date_to")),
		GroupBy:            strings.TrimSpace(values.Get("group_by")),
		PaymentType:        strings.TrimSpace(values.Get("payment_type")),
		PaymentStatus:      strings.TrimSpace(values.Get("payment_status")),
		OverdueOnly:        overdueOnly,
		RequesterUserID:    authUser.UserID,
		IncludeAllProjects: authUser.Role == core.RoleAdmin,
	}
	if input.IncludeAllProjects {
		input.RequesterUserID = ""
	}
	return input, nil
}
