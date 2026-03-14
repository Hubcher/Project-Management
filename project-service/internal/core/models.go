package core

import "time"

type Project struct {
	ContractNumber int64
	Name           string
	StartDate      time.Time
	Deadline       time.Time
	Price          string
	UserID         string
	CreatedAt      time.Time
}

type CreateProjectInput struct {
	Name      string
	StartDate time.Time
	Deadline  time.Time
	Price     string
	UserID    string
}

type UpdateProjectInput struct {
	ContractNumber int64
	Name           string
	StartDate      time.Time
	Deadline       time.Time
	Price          string
	UserID         string
}
