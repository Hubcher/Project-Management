package core

import "context"

type ProjectRepo interface {
	CreateProject(ctx context.Context, dto NewProjectDto) (*Project, error)

	GetAllProjects(ctx context.Context) ([]Project, error)

	GetProjectById(ctx context.Context, id int) (*Project, error)

	GetProjectByName(ctx context.Context, name string) (*Project, error)

	UpdateProject(ctx context.Context, project *Project) (*Project, error)

	DeleteProject(ctx context.Context, id int) (int, error)
}
