package core

import "context"

type DB interface {
	CreateUser(ctx context.Context, input CreateUserInput) (*User, error)
	GetUser(ctx context.Context, id string) (*User, error)
	ListUsers(ctx context.Context, role string) ([]User, error)
	UpdateUser(ctx context.Context, input UpdateUserInput) (*User, error)
	DeleteUser(ctx context.Context, id string) error
}

type UserService interface {
	CreateUser(ctx context.Context, input CreateUserInput) (*User, error)
	GetUser(ctx context.Context, id string) (*User, error)
	ListUsers(ctx context.Context, role string) ([]User, error)
	UpdateUser(ctx context.Context, input UpdateUserInput) (*User, error)
	DeleteUser(ctx context.Context, id string) error
}
