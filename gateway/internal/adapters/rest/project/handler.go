package project

import (
	"context"
	"net/http"
	"strings"

	"github.com/Hubcher/project-management/gateway/internal/adapters/rest/httpx"
	"github.com/Hubcher/project-management/gateway/internal/adapters/rest/middleware"
	"github.com/Hubcher/project-management/gateway/internal/core"
)

// NewCreateProjectHandler godoc
// @Summary Create project
// @Description Creates a project. Administrators may assign any manager, regular users may create projects only for themselves.
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateProjectRequest true "Project payload"
// @Success 201 {object} ProjectResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/projects [post]
func NewCreateProjectHandler(projects core.ProjectDirectory, users core.UserDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		var req CreateProjectRequest
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		managerID, err := resolveProjectManager(authUser, req.ManagerID)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if err = ensureKnownUser(r.Context(), users, managerID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		project, err := projects.CreateProject(r.Context(), toCreateProjectInput(req, managerID))
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusCreated, toProjectResponse(project))
	}
}

// NewGetProjectHandler godoc
// @Summary Get project
// @Description Returns a project by ID. Accessible to administrators, the project manager, and active project members.
// @Tags projects
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Success 200 {object} ProjectResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/projects/{id} [get]
func NewGetProjectHandler(projects core.ProjectDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectID := strings.TrimSpace(r.PathValue("id"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		project, err := loadReadableProject(r.Context(), projects, authUser, projectID)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, toProjectResponse(project))
	}
}

// NewListProjectsHandler godoc
// @Summary List projects
// @Description Returns projects visible to the caller. Administrators may use the participant filter, regular users see only their own related projects.
// @Tags projects
// @Produce json
// @Security BearerAuth
// @Param participant_user_id query string false "Participant user ID filter"
// @Success 200 {object} ListProjectsResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/projects [get]
func NewListProjectsHandler(projects core.ProjectDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		participantUserID, err := resolveParticipantFilter(authUser, r.URL.Query().Get("participant_user_id"))
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		list, err := projects.ListProjects(r.Context(), participantUserID)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, toListProjectsResponse(list))
	}
}

// NewUpdateProjectHandler godoc
// @Summary Update project
// @Description Updates a project. Accessible to administrators and the current project manager.
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Param request body UpdateProjectRequest true "Project update payload"
// @Success 200 {object} ProjectResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/projects/{id} [put]
func NewUpdateProjectHandler(projects core.ProjectDirectory, users core.UserDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectID := strings.TrimSpace(r.PathValue("id"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		project, err := loadWritableProject(r.Context(), projects, authUser, projectID)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		var req UpdateProjectRequest
		if err = httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		managerID, err := resolveProjectManager(authUser, req.ManagerID)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if authUser.Role != core.RoleAdmin && managerID != project.ManagerID {
			httpx.WriteError(w, http.StatusForbidden, "access denied")
			return
		}
		if err = ensureKnownUser(r.Context(), users, managerID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		updated, err := projects.UpdateProject(r.Context(), toUpdateProjectInput(projectID, req, managerID))
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, toProjectResponse(updated))
	}
}

// NewDeleteProjectHandler godoc
// @Summary Delete project
// @Description Deletes a project. Accessible to administrators and the current project manager.
// @Tags projects
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Success 204 {string} string "No Content"
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/projects/{id} [delete]
func NewDeleteProjectHandler(projects core.ProjectDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectID := strings.TrimSpace(r.PathValue("id"))
		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		if _, err := loadWritableProject(r.Context(), projects, authUser, projectID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		if err := projects.DeleteProject(r.Context(), projectID); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func loadReadableProject(ctx context.Context, projects core.ProjectDirectory, authUser core.AuthUser, projectID string) (*core.Project, error) {
	project, err := projects.GetProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if authUser.Role == core.RoleAdmin || project.ManagerID == authUser.UserID {
		return project, nil
	}
	members, err := projects.ListMembers(ctx, project.ID)
	if err != nil {
		return nil, err
	}
	if isActiveProjectMember(members, authUser.UserID) {
		return project, nil
	}
	return nil, core.NewStatusError(http.StatusForbidden, "access denied")
}

func loadWritableProject(ctx context.Context, projects core.ProjectDirectory, authUser core.AuthUser, projectID string) (*core.Project, error) {
	project, err := projects.GetProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if authUser.Role == core.RoleAdmin || project.ManagerID == authUser.UserID {
		return project, nil
	}
	return nil, core.NewStatusError(http.StatusForbidden, "access denied")
}

func resolveProjectManager(authUser core.AuthUser, requestedManagerID string) (string, error) {
	requestedManagerID = strings.TrimSpace(requestedManagerID)
	if authUser.Role == core.RoleAdmin {
		if requestedManagerID == "" {
			return authUser.UserID, nil
		}
		return requestedManagerID, nil
	}
	if requestedManagerID != "" && requestedManagerID != authUser.UserID {
		return "", core.NewStatusError(http.StatusForbidden, "access denied")
	}
	return authUser.UserID, nil
}

func resolveParticipantFilter(authUser core.AuthUser, requested string) (string, error) {
	requested = strings.TrimSpace(requested)
	if authUser.Role == core.RoleAdmin {
		return requested, nil
	}
	if requested != "" && requested != authUser.UserID {
		return "", core.NewStatusError(http.StatusForbidden, "access denied")
	}
	return authUser.UserID, nil
}

func ensureKnownUser(ctx context.Context, users core.UserDirectory, userID string) error {
	_, err := users.GetUser(ctx, strings.TrimSpace(userID))
	return err
}

func isActiveProjectMember(members []core.ProjectMember, userID string) bool {
	for _, member := range members {
		if member.UserID == userID && member.IsActive {
			return true
		}
	}
	return false
}