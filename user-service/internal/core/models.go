package core

import "time"

type User struct {
	ID        string
	Name      string
	Email     string
	Password  string
	Role      string
	CreatedAt time.Time
}

type CreateUserInput struct {
	ID       string
	Name     string
	Email    string
	Password string
	Role     string
}

type UpdateUserInput struct {
	ID       string
	Name     string
	Email    string
	Password string
	Role     string
}
