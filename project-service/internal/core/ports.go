package core

import "context"

type ProjectRepository interface {
	CreateProject(ctx context.Context, input CreateProjectInput) (*Project, error)
	GetProject(ctx context.Context, contractNumber int64) (*Project, error)
	ListProjects(ctx context.Context, userID string) ([]Project, error)
	UpdateProject(ctx context.Context, input UpdateProjectInput) (*Project, error)
	DeleteProject(ctx context.Context, contractNumber int64) error
}

type ProjectService interface {
	CreateProject(ctx context.Context, input CreateProjectInput) (*Project, error)
	GetProject(ctx context.Context, contractNumber int64) (*Project, error)
	ListProjects(ctx context.Context, userID string) ([]Project, error)
	UpdateProject(ctx context.Context, input UpdateProjectInput) (*Project, error)
	DeleteProject(ctx context.Context, contractNumber int64) error
}
