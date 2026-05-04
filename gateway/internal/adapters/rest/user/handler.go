package user

import (
	"net/http"

	"github.com/Hubcher/project-management/gateway/internal/adapters/rest/httpx"
	"github.com/Hubcher/project-management/gateway/internal/adapters/rest/middleware"
	"github.com/Hubcher/project-management/gateway/internal/core"
)

// NewCreateUserHandler godoc
// @Summary Create a managed user
// @Description Creates a new user, manager or administrator account. This endpoint is available only to administrators.
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateUserRequest true "Managed user payload"
// @Success 201 {object} ManagedUserResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/users [post]
func NewCreateUserHandler(identity *core.IdentityService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateUserRequest
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := identity.CreateManagedUser(r.Context(), toRegisterInput(req))
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		httpx.WriteJSON(w, http.StatusCreated, toManagedUserResponse(result))
	}
}

// NewGetUserByIDHandler godoc
// @Summary Get user profile
// @Description Returns a user profile by ID. Administrators can read any profile, regular users can read only their own profile.
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} UserResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/users/{id} [get]
func NewGetUserByIDHandler(users core.UserDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			httpx.WriteError(w, http.StatusBadRequest, "invalid user id")
			return
		}

		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}
		if authUser.Role != core.RoleAdmin && authUser.UserID != id {
			httpx.WriteError(w, http.StatusForbidden, "access denied")
			return
		}

		user, err := users.GetUser(r.Context(), id)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		httpx.WriteJSON(w, http.StatusOK, toUserResponse(user))
	}
}

// NewListUsersHandler godoc
// @Summary List user profiles
// @Description Returns all user profiles. This endpoint is available only to administrators.
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ListUsersResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/users [get]
func NewListUsersHandler(users core.UserDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profiles, err := users.ListUsers(r.Context())
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		httpx.WriteJSON(w, http.StatusOK, toListUsersResponse(profiles))
	}
}

// NewUpdateUserHandler godoc
// @Summary Update user profile
// @Description Updates a user profile. Administrators can update any profile, regular users can update only their own profile.
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body UpdateUserRequest true "Profile update payload"
// @Success 200 {object} UserResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/users/{id} [put]
func NewUpdateUserHandler(users core.UserDirectory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			httpx.WriteError(w, http.StatusBadRequest, "invalid user id")
			return
		}

		authUser, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}
		if authUser.Role != core.RoleAdmin && authUser.UserID != id {
			httpx.WriteError(w, http.StatusForbidden, "access denied")
			return
		}

		var req UpdateUserRequest
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		user, err := users.UpdateUser(r.Context(), toUpdateInput(id, req))
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		httpx.WriteJSON(w, http.StatusOK, toUserResponse(user))
	}
}

// NewDeleteUserHandler godoc
// @Summary Delete user
// @Description Deletes a user in both auth-service and user-service. This endpoint is available only to administrators.
// @Tags users
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/users/{id} [delete]
func NewDeleteUserHandler(identity *core.IdentityService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			httpx.WriteError(w, http.StatusBadRequest, "invalid user id")
			return
		}

		if err := identity.DeleteUser(r.Context(), id); err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
