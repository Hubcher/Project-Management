package project

import (
	"net/http"
	"strings"

	"github.com/Hubcher/project-management/gateway/internal/adapters/rest/httpx"
	"github.com/Hubcher/project-management/gateway/internal/adapters/rest/middleware"
)

// NewCreateEventHandler godoc
// @Summary Create project event
// @Description Creates a project event. Accessible to administrators and the project manager.
// @Tags project-events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param projectId path string true "Project ID"
// @Param request body CreateProjectEventRequest true "Project event payload"
// @Success 201 {object} ProjectEventResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/projects/{projectId}/events [post]
func NewCreateEventHandler(projects ProjectHandlerService) http.HandlerFunc {
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

		var req CreateProjectEventRequest
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		projectEvent, err := projects.CreateEvent(r.Context(), toCreateEventInput(projectID, req))
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusCreated, toEventResponse(projectEvent))
	}
}

// NewListEventsHandler godoc
// @Summary List project events
// @Description Returns project events. Accessible to administrators, the project manager, and active project members.
// @Tags project-events
// @Produce json
// @Security BearerAuth
// @Param projectId path string true "Project ID"
// @Success 200 {object} ListProjectEventsResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/projects/{projectId}/events [get]
func NewListEventsHandler(projects ProjectHandlerService) http.HandlerFunc {
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

		events, err := projects.ListEvents(r.Context(), projectID)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, toListEventsResponse(events))
	}
}

// NewGetEventHandler godoc
// @Summary Get project event
// @Description Returns a project event by ID. Accessible to administrators, the project manager, and active project members.
// @Tags project-events
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project event ID"
// @Success 200 {object} ProjectEventResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/project-events/{id} [get]
func NewGetEventHandler(projects ProjectHandlerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.PathValue("id"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		projectEvent, err := projects.GetEvent(r.Context(), id)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if _, err = loadReadableProject(r.Context(), projects, authUser, projectEvent.ProjectID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, toEventResponse(projectEvent))
	}
}

// NewUpdateEventHandler godoc
// @Summary Update project event
// @Description Updates a project event. Accessible to administrators and the project manager.
// @Tags project-events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project event ID"
// @Param request body UpdateProjectEventRequest true "Project event update payload"
// @Success 200 {object} ProjectEventResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/project-events/{id} [put]
func NewUpdateEventHandler(projects ProjectHandlerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.PathValue("id"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		projectEvent, err := projects.GetEvent(r.Context(), id)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if _, err = loadWritableProject(r.Context(), projects, authUser, projectEvent.ProjectID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		var req UpdateProjectEventRequest
		if err = httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		updated, err := projects.UpdateEvent(r.Context(), toUpdateEventInput(id, req))
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, toEventResponse(updated))
	}
}

// NewDeleteEventHandler godoc
// @Summary Delete project event
// @Description Deletes a project event. Accessible to administrators and the project manager.
// @Tags project-events
// @Security BearerAuth
// @Param id path string true "Project event ID"
// @Success 204 {string} string "No Content"
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/project-events/{id} [delete]
func NewDeleteEventHandler(projects ProjectHandlerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.PathValue("id"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		projectEvent, err := projects.GetEvent(r.Context(), id)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if _, err = loadWritableProject(r.Context(), projects, authUser, projectEvent.ProjectID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if err = projects.DeleteEvent(r.Context(), id); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}