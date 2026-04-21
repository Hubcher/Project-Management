package core

import "time"

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type Account struct {
	UserID            string
	Email             string
	PasswordHash      string
	Role              Role
	IsActive          bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
	PasswordChangedAt time.Time
}

type Claims struct {
	UserID string
	Role   Role
	Email  string
}
