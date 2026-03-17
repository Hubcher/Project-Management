package core

import (
	"context"
	"time"
)

type AuthRepository interface {
	CreateAccount(ctx context.Context, account AuthAccount) error
	GetAccountByEmail(ctx context.Context, email string) (AuthAccount, error)
	GetAccountByUserID(ctx context.Context, userID string) (AuthAccount, error)
	DeleteAccount(ctx context.Context, userID string) error

	CreateRefreshSession(ctx context.Context, session RefreshSession) error
	GetRefreshSessionByTokenHash(ctx context.Context, tokenHash string) (RefreshSession, error)
	RevokeRefreshSession(ctx context.Context, sessionID string, revokedAt time.Time) error
	RevokeAllRefreshSessionsByUserID(ctx context.Context, userID string, revokedAt time.Time) error
}

type UserProvisioner interface {
	CreateUser(ctx context.Context, id string, name string) error
}
