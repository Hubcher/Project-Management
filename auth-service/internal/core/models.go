package core

import "time"

type AuthAccount struct {
	UserID            string    `db:"user_id"`
	Email             string    `db:"email"`
	PasswordHash      string    `db:"password_hash"`
	Role              string    `db:"role"`
	IsActive          bool      `db:"is_active"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
	PasswordChangedAt time.Time `db:"password_changed_at"`
}

type RefreshSession struct {
	ID        string     `db:"id"`
	UserID    string     `db:"user_id"`
	TokenHash string     `db:"token_hash"`
	UserAgent *string    `db:"user_agent"`
	IP        *string    `db:"ip"`
	ExpiresAt time.Time  `db:"expires_at"`
	CreatedAt time.Time  `db:"created_at"`
	RevokedAt *time.Time `db:"revoked_at"`
}
