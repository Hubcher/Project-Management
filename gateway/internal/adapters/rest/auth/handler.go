package auth

import (
	"net/http"

	"github.com/Hubcher/project-management/gateway/internal/adapters/rest/httpx"
	"github.com/Hubcher/project-management/gateway/internal/adapters/rest/middleware"
	"github.com/Hubcher/project-management/gateway/internal/core"
)

// NewRegisterHandler godoc
// @Summary Register a user
// @Description Creates a regular user account and returns a JWT access token.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration payload"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/auth/register [post]
func NewRegisterHandler(service *core.IdentityService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := service.Register(r.Context(), toRegisterInput(req))
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		httpx.WriteJSON(w, http.StatusCreated, toRegisterResponse(result))
	}
}

// NewLoginHandler godoc
// @Summary Login
// @Description Authenticates a user by email and password and returns a JWT access token.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login payload"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 502 {object} httpx.ErrorResponse
// @Router /api/auth/login [post]
func NewLoginHandler(service *core.IdentityService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := service.Login(r.Context(), req.Email, req.Password)
		if err != nil {
			httpx.WriteAnyError(w, err)
			return
		}

		httpx.WriteJSON(w, http.StatusOK, toLoginResponse(result))
	}
}

// NewMeHandler godoc
// @Summary Get current auth user
// @Description Returns JWT claims of the currently authenticated user.
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} AuthUserResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /api/auth/me [get]
func NewMeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := middleware.CurrentUser(r.Context())
		if !ok {
			httpx.WriteError(w, http.StatusUnauthorized, "missing auth context")
			return
		}

		httpx.WriteJSON(w, http.StatusOK, toAuthUserResponse(user))
	}
}