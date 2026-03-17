package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
)

type ctxKey string

const claimsKey ctxKey = "auth_claims"

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type AuthMiddleware struct {
	log    *slog.Logger
	issuer string
	secret []byte
}

func NewAuthMiddleware(log *slog.Logger, issuer, secret string) *AuthMiddleware {
	return &AuthMiddleware{
		log:    log,
		issuer: issuer,
		secret: []byte(secret),
	}
}

func Chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func (m *AuthMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawToken := bearerToken(r.Header.Get("Authorization"))
		if rawToken == "" {
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (any, error) {
			if token.Method == nil || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return m.secret, nil
		}, jwt.WithIssuer(m.issuer))
		if err != nil || !token.Valid {
			m.log.Warn("token validation failed", "error", err)
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		if claims.Subject == "" {
			http.Error(w, "token subject is empty", http.StatusUnauthorized)
			return
		}
		if claims.Role == "" {
			http.Error(w, "token role is empty", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, *claims)

		// Эти metadata автоматически уйдут в gRPC, если handler потом вызовет client.Method(r.Context(), ...)
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(
			"x-user-id", claims.Subject,
			"x-user-role", claims.Role,
		))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequireRoles(allowedRoles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(allowedRoles))
	for _, role := range allowedRoles {
		allowed[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetClaims(r.Context())
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			if _, exists := allowed[claims.Role]; !exists {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetClaims(ctx context.Context) (Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(Claims)
	return claims, ok
}

func bearerToken(header string) string {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}
