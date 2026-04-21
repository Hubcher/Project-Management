package middleware

import (
    "context"
    "encoding/json"
    "net/http"
    "strings"

    "github.com/Hubcher/project-management/gateway/internal/core"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type ctxKey string

const authUserKey ctxKey = "auth_user"

type AuthMiddleware struct {
    validator core.AuthClient
}

func NewAuthMiddleware(validator core.AuthClient) *AuthMiddleware {
    return &AuthMiddleware{validator: validator}
}

func (m *AuthMiddleware) Auth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        header := strings.TrimSpace(r.Header.Get("Authorization"))
        if !strings.HasPrefix(strings.ToLower(header), "bearer ") {
            writeError(w, http.StatusUnauthorized, "missing bearer token")
            return
        }

        token := strings.TrimSpace(header[len("Bearer "):])
        if token == "" {
            writeError(w, http.StatusUnauthorized, "missing bearer token")
            return
        }

        user, err := m.validator.Validate(r.Context(), token)
        if err != nil {
            switch status.Code(err) {
            case codes.Unauthenticated:
                writeError(w, http.StatusUnauthorized, "invalid token")
            case codes.PermissionDenied:
                writeError(w, http.StatusForbidden, "access denied")
            default:
                writeError(w, http.StatusBadGateway, "auth-service is unavailable")
            }
            return
        }

        next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), authUserKey, user)))
    })
}

func (m *AuthMiddleware) RequireRoles(roles ...core.Role) func(http.Handler) http.Handler {
    allowed := make(map[core.Role]struct{}, len(roles))
    for _, role := range roles {
        allowed[role] = struct{}{}
    }

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user, ok := CurrentUser(r.Context())
            if !ok {
                writeError(w, http.StatusUnauthorized, "missing auth context")
                return
            }
            if _, exists := allowed[user.Role]; !exists {
                writeError(w, http.StatusForbidden, "access denied")
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

func CurrentUser(ctx context.Context) (core.AuthUser, bool) {
    user, ok := ctx.Value(authUserKey).(core.AuthUser)
    return user, ok
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    _ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
