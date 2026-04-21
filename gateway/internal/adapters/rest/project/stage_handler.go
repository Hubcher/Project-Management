package project

import (
	"net/http"
	"strings"

	"github.com/Hubcher/project-management/gateway/internal/adapters/rest/httpx"
	"github.com/Hubcher/project-management/gateway/internal/adapters/rest/middleware"
)

// NewCreateStageHandler godoc
// @Summary Create project stage
// @Description Creates a stage inside a project. Accessible to administrators and the project manager.
// @Tags project-stages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param projectId path string true "Project ID"
// @Param request body CreateProjectStageRequest true "Project stage payload"
// @Success 201 {object} ProjectStageResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/projects/{projectId}/stages [post]
func NewCreateStageHandler(projects ProjectHandlerService) http.HandlerFunc {
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

		var req CreateProjectStageRequest
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		stage, err := projects.CreateStage(r.Context(), toCreateStageInput(projectID, req))
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusCreated, toStageResponse(stage))
	}
}

// NewListStagesHandler godoc
// @Summary List project stages
// @Description Returns all stages of a project. Accessible to administrators, the project manager, and active project members.
// @Tags project-stages
// @Produce json
// @Security BearerAuth
// @Param projectId path string true "Project ID"
// @Success 200 {object} ListProjectStagesResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/projects/{projectId}/stages [get]
func NewListStagesHandler(projects ProjectHandlerService) http.HandlerFunc {
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

		stages, err := projects.ListStages(r.Context(), projectID)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, toListStagesResponse(stages))
	}
}

// NewGetStageHandler godoc
// @Summary Get project stage
// @Description Returns a project stage by ID. Accessible to administrators, the project manager, and active project members.
// @Tags project-stages
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project stage ID"
// @Success 200 {object} ProjectStageResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/project-stages/{id} [get]
func NewGetStageHandler(projects ProjectHandlerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.PathValue("id"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		stage, err := projects.GetStage(r.Context(), id)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if _, err = loadReadableProject(r.Context(), projects, authUser, stage.ProjectID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, toStageResponse(stage))
	}
}

// NewUpdateStageHandler godoc
// @Summary Update project stage
// @Description Updates a project stage. Accessible to administrators and the project manager.
// @Tags project-stages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project stage ID"
// @Param request body UpdateProjectStageRequest true "Project stage update payload"
// @Success 200 {object} ProjectStageResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/project-stages/{id} [put]
func NewUpdateStageHandler(projects ProjectHandlerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.PathValue("id"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		stage, err := projects.GetStage(r.Context(), id)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if _, err = loadWritableProject(r.Context(), projects, authUser, stage.ProjectID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		var req UpdateProjectStageRequest
		if err = httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		updated, err := projects.UpdateStage(r.Context(), toUpdateStageInput(id, req))
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, toStageResponse(updated))
	}
}

// NewDeleteStageHandler godoc
// @Summary Delete project stage
// @Description Deletes a project stage. Accessible to administrators and the project manager.
// @Tags project-stages
// @Security BearerAuth
// @Param id path string true "Project stage ID"
// @Success 204 {string} string "No Content"
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/project-stages/{id} [delete]
func NewDeleteStageHandler(projects ProjectHandlerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.PathValue("id"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		stage, err := projects.GetStage(r.Context(), id)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if _, err = loadWritableProject(r.Context(), projects, authUser, stage.ProjectID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if err = projects.DeleteStage(r.Context(), id); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}