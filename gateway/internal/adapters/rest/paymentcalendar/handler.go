package paymentcalendar

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/Hubcher/project-management/gateway/internal/adapters/rest/httpx"
	"github.com/Hubcher/project-management/gateway/internal/adapters/rest/middleware"
	"github.com/Hubcher/project-management/gateway/internal/core"
)

type queryValues struct {
	stageID     string
	paymentType string
	status      string
	dateFrom    string
	dateTo      string
	overdueOnly bool
}

// NewCreatePaymentHandler godoc
// @Summary Create project payment
// @Description Creates a planned income or expense for a project. Accessible to administrators and the project manager.
// @Tags payment-calendar
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param projectId path string true "Project ID"
// @Param request body CreatePaymentRequest true "Payment payload"
// @Success 201 {object} PaymentResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/projects/{projectId}/payments [post]
func NewCreatePaymentHandler(payments PaymentHandlerService, projects core.ProjectDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectID := strings.TrimSpace(r.PathValue("projectId"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}
		if _, err := loadWritableProject(r.Context(), projects, authUser, projectID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		var req CreatePaymentRequest
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		payment, err := payments.CreatePayment(r.Context(), toCreatePaymentInput(projectID, authUser.UserID, req))
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusCreated, toPaymentResponse(payment))
	}
}

// NewListPaymentsHandler godoc
// @Summary List project payments
// @Description Returns project payment calendar entries visible to the caller.
// @Tags payment-calendar
// @Produce json
// @Security BearerAuth
// @Param projectId path string true "Project ID"
// @Param stage_id query string false "Project stage ID filter"
// @Param type query string false "Payment type: income or expense"
// @Param status query string false "Stored payment status: planned, paid or cancelled"
// @Param date_from query string false "Lower planned date bound in YYYY-MM-DD format"
// @Param date_to query string false "Upper planned date bound in YYYY-MM-DD format"
// @Param overdue_only query bool false "Return only computed overdue payments"
// @Success 200 {object} ListPaymentsResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/projects/{projectId}/payments [get]
func NewListPaymentsHandler(payments PaymentHandlerService, projects core.ProjectDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectID := strings.TrimSpace(r.PathValue("projectId"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}
		if _, err := loadReadableProject(r.Context(), projects, authUser, projectID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		query, err := parseQueryValues(r)
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		list, err := payments.ListPayments(r.Context(), toListPaymentsInput(projectID, query))
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, toListPaymentsResponse(list))
	}
}

// NewGetPaymentHandler godoc
// @Summary Get project payment
// @Description Returns a payment by ID. Accessible to administrators, the project manager, and active project members.
// @Tags payment-calendar
// @Produce json
// @Security BearerAuth
// @Param id path string true "Payment ID"
// @Success 200 {object} PaymentResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/payments/{id} [get]
func NewGetPaymentHandler(payments PaymentHandlerService, projects core.ProjectDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paymentID := strings.TrimSpace(r.PathValue("id"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		payment, err := payments.GetPayment(r.Context(), paymentID)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if _, err = loadReadableProject(r.Context(), projects, authUser, payment.ProjectID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, toPaymentResponse(payment))
	}
}

// NewUpdatePaymentHandler godoc
// @Summary Update project payment
// @Description Updates a payment. Accessible to administrators and the project manager.
// @Tags payment-calendar
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Payment ID"
// @Param request body UpdatePaymentRequest true "Payment update payload"
// @Success 200 {object} PaymentResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/payments/{id} [put]
func NewUpdatePaymentHandler(payments PaymentHandlerService, projects core.ProjectDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paymentID := strings.TrimSpace(r.PathValue("id"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		current, err := payments.GetPayment(r.Context(), paymentID)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if _, err = loadWritableProject(r.Context(), projects, authUser, current.ProjectID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		var req UpdatePaymentRequest
		if err = httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		if req.Status != nil && current.Status != "paid" && strings.EqualFold(strings.TrimSpace(*req.Status), "paid") {
			httpx.WriteError(w, http.StatusBadRequest, "use pay endpoint to mark payment paid")
			return
		}

		updated, err := payments.UpdatePayment(r.Context(), toUpdatePaymentInput(paymentID, current, req))
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, toPaymentResponse(updated))
	}
}

// NewDeletePaymentHandler godoc
// @Summary Delete project payment
// @Description Deletes a payment. Accessible to administrators and the project manager.
// @Tags payment-calendar
// @Security BearerAuth
// @Param id path string true "Payment ID"
// @Success 204 {string} string "No Content"
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/payments/{id} [delete]
func NewDeletePaymentHandler(payments PaymentHandlerService, projects core.ProjectDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paymentID := strings.TrimSpace(r.PathValue("id"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		payment, err := payments.GetPayment(r.Context(), paymentID)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if _, err = loadWritableProject(r.Context(), projects, authUser, payment.ProjectID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if err = payments.DeletePayment(r.Context(), paymentID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// NewMarkPaymentPaidHandler godoc
// @Summary Mark payment paid
// @Description Marks a planned payment as paid and stores the actual payment date.
// @Tags payment-calendar
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Payment ID"
// @Param request body MarkPaymentPaidRequest false "Actual payment date payload"
// @Success 200 {object} PaymentResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/payments/{id}/pay [post]
func NewMarkPaymentPaidHandler(payments PaymentHandlerService, projects core.ProjectDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paymentID := strings.TrimSpace(r.PathValue("id"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		payment, err := payments.GetPayment(r.Context(), paymentID)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if _, err = loadWritableProject(r.Context(), projects, authUser, payment.ProjectID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		var req MarkPaymentPaidRequest
		if r.ContentLength != 0 {
			if err = httpx.DecodeJSON(r, &req); err != nil {
				httpx.WriteError(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		paid, err := payments.MarkPaymentPaid(r.Context(), core.MarkPaymentPaidInput{ID: paymentID, ActualDate: req.ActualDate, PaidBy: authUser.UserID})
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, toPaymentResponse(paid))
	}
}

// NewProjectSummaryHandler godoc
// @Summary Get project financial summary
// @Description Returns planned, paid and computed overdue totals for a project.
// @Tags payment-calendar
// @Produce json
// @Security BearerAuth
// @Param projectId path string true "Project ID"
// @Success 200 {object} ProjectFinancialSummaryResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/projects/{projectId}/payment-summary [get]
func NewProjectSummaryHandler(payments PaymentHandlerService, projects core.ProjectDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectID := strings.TrimSpace(r.PathValue("projectId"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}
		if _, err := loadReadableProject(r.Context(), projects, authUser, projectID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		summary, err := payments.GetProjectSummary(r.Context(), projectID)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, toSummaryResponse(summary))
	}
}

func parseQueryValues(r *http.Request) (queryValues, error) {
	values := r.URL.Query()
	overdueOnly := false
	if raw := strings.TrimSpace(values.Get("overdue_only")); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			return queryValues{}, err
		}
		overdueOnly = parsed
	}
	return queryValues{
		stageID:     values.Get("stage_id"),
		paymentType: values.Get("type"),
		status:      values.Get("status"),
		dateFrom:    values.Get("date_from"),
		dateTo:      values.Get("date_to"),
		overdueOnly: overdueOnly,
	}, nil
}

func loadReadableProject(ctx context.Context, projects core.ProjectDirectory, authUser core.AuthUser, projectID string) (*core.Project, error) {
	project, err := projects.GetProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if authUser.Role == core.RoleAdmin || isProjectManager(authUser, project) {
		return project, nil
	}
	members, err := projects.ListMembers(ctx, project.ID)
	if err != nil {
		return nil, err
	}
	for _, member := range members {
		if member.UserID == authUser.UserID && member.IsActive {
			return project, nil
		}
	}
	return nil, core.NewStatusError(http.StatusForbidden, "access denied")
}

func loadWritableProject(ctx context.Context, projects core.ProjectDirectory, authUser core.AuthUser, projectID string) (*core.Project, error) {
	project, err := projects.GetProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if authUser.Role == core.RoleAdmin || isProjectManager(authUser, project) {
		return project, nil
	}
	return nil, core.NewStatusError(http.StatusForbidden, "access denied")
}

func isProjectManager(authUser core.AuthUser, project *core.Project) bool {
	return authUser.Role == core.RoleManager && project.ManagerID == authUser.UserID
}
