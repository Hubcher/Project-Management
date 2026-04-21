package project

import (
	"net/http"
	"strings"

	"github.com/Hubcher/project-management/gateway/internal/adapters/rest/httpx"
	"github.com/Hubcher/project-management/gateway/internal/adapters/rest/middleware"
	"github.com/Hubcher/project-management/gateway/internal/core"
)

// NewCreateMemberHandler godoc
// @Summary Add project member
// @Description Adds a user to a project team. Accessible to administrators and the project manager.
// @Tags project-members
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param projectId path string true "Project ID"
// @Param request body CreateProjectMemberRequest true "Project member payload"
// @Success 201 {object} ProjectMemberResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/projects/{projectId}/members [post]
func NewCreateMemberHandler(projects ProjectHandlerService, users core.UserDirectory) http.HandlerFunc {
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

		var req CreateProjectMemberRequest
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := ensureKnownUser(r.Context(), users, req.UserID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		member, err := projects.CreateMember(r.Context(), toCreateMemberInput(projectID, req))
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusCreated, toMemberResponse(member))
	}
}

// NewListMembersHandler godoc
// @Summary List project members
// @Description Returns project members. Accessible to administrators, the project manager, and active project members.
// @Tags project-members
// @Produce json
// @Security BearerAuth
// @Param projectId path string true "Project ID"
// @Success 200 {object} ListProjectMembersResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/projects/{projectId}/members [get]
func NewListMembersHandler(projects ProjectHandlerService) http.HandlerFunc {
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

		members, err := projects.ListMembers(r.Context(), projectID)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, toListMembersResponse(members))
	}
}

// NewGetMemberHandler godoc
// @Summary Get project member
// @Description Returns a project member by ID. Accessible to administrators, the project manager, and active project members.
// @Tags project-members
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project member ID"
// @Success 200 {object} ProjectMemberResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/project-members/{id} [get]
func NewGetMemberHandler(projects ProjectHandlerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.PathValue("id"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		member, err := projects.GetMember(r.Context(), id)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if _, err = loadReadableProject(r.Context(), projects, authUser, member.ProjectID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, toMemberResponse(member))
	}
}

// NewUpdateMemberHandler godoc
// @Summary Update project member
// @Description Updates project member attributes. Accessible to administrators and the project manager.
// @Tags project-members
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project member ID"
// @Param request body UpdateProjectMemberRequest true "Project member update payload"
// @Success 200 {object} ProjectMemberResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/project-members/{id} [put]
func NewUpdateMemberHandler(projects ProjectHandlerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.PathValue("id"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		member, err := projects.GetMember(r.Context(), id)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if _, err = loadWritableProject(r.Context(), projects, authUser, member.ProjectID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		var req UpdateProjectMemberRequest
		if err = httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		updated, err := projects.UpdateMember(r.Context(), toUpdateMemberInput(id, req))
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, toMemberResponse(updated))
	}
}

// NewDeleteMemberHandler godoc
// @Summary Delete project member
// @Description Deletes a project member. Accessible to administrators and the project manager.
// @Tags project-members
// @Security BearerAuth
// @Param id path string true "Project member ID"
// @Success 204 {string} string "No Content"
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/project-members/{id} [delete]
func NewDeleteMemberHandler(projects ProjectHandlerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimSpace(r.PathValue("id"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		member, err := projects.GetMember(r.Context(), id)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if _, err = loadWritableProject(r.Context(), projects, authUser, member.ProjectID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if err = projects.DeleteMember(r.Context(), id); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}