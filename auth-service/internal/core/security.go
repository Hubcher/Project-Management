package core

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordManager struct {
	Cost int
}

func (m BcryptPasswordManager) Hash(password string) (string, error) {
	cost := m.Cost
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (m BcryptPasswordManager) Compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

type JWTManager struct {
	Secret   string
	Issuer   string
	Audience string
	TTL      time.Duration
}

type jwtClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func (m JWTManager) Generate(claims Claims) (string, error) {
	ttl := m.TTL
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims{
		UserID: claims.UserID,
		Role:   string(claims.Role),
		Email:  claims.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.Issuer,
			Audience:  jwt.ClaimStrings{m.Audience},
			Subject:   claims.UserID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	})
	return token.SignedString([]byte(m.Secret))
}

func (m JWTManager) Parse(token string) (Claims, error) {
	parsed, err := jwt.ParseWithClaims(token, &jwtClaims{}, func(_ *jwt.Token) (any, error) {
		return []byte(m.Secret), nil
	}, jwt.WithIssuer(m.Issuer), jwt.WithAudience(m.Audience))
	if err != nil {
		return Claims{}, ErrInvalidToken
	}

	claims, ok := parsed.Claims.(*jwtClaims)
	if !ok || !parsed.Valid {
		return Claims{}, ErrInvalidToken
	}
	if claims.UserID == "" || claims.Email == "" || claims.Role == "" {
		return Claims{}, ErrInvalidToken
	}

	role := Role(claims.Role)
	if role != RoleAdmin && role != RoleUser {
		return Claims{}, fmt.Errorf("%w: unknown role", ErrInvalidToken)
	}

	return Claims{
		UserID: claims.UserID,
		Role:   role,
		Email:  claims.Email,
	}, nil
}
