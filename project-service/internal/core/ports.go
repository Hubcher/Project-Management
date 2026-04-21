package core

import "context"

type ProjectRepository interface {
	CreateProject(ctx context.Context, input CreateProjectInput) (*Project, error)
	GetProject(ctx context.Context, id string) (*Project, error)
	ListProjects(ctx context.Context, participantUserID string) ([]Project, error)
	UpdateProject(ctx context.Context, input UpdateProjectInput) (*Project, error)
	DeleteProject(ctx context.Context, id string) error

	CreateStage(ctx context.Context, input CreateProjectStageInput) (*ProjectStage, error)
	GetStage(ctx context.Context, id string) (*ProjectStage, error)
	ListStages(ctx context.Context, projectID string) ([]ProjectStage, error)
	UpdateStage(ctx context.Context, input UpdateProjectStageInput) (*ProjectStage, error)
	DeleteStage(ctx context.Context, id string) error

	CreateMember(ctx context.Context, input CreateProjectMemberInput) (*ProjectMember, error)
	GetMember(ctx context.Context, id string) (*ProjectMember, error)
	ListMembers(ctx context.Context, projectID string) ([]ProjectMember, error)
	UpdateMember(ctx context.Context, input UpdateProjectMemberInput) (*ProjectMember, error)
	DeleteMember(ctx context.Context, id string) error

	CreateEvent(ctx context.Context, input CreateProjectEventInput) (*ProjectEvent, error)
	GetEvent(ctx context.Context, id string) (*ProjectEvent, error)
	ListEvents(ctx context.Context, projectID string) ([]ProjectEvent, error)
	UpdateEvent(ctx context.Context, input UpdateProjectEventInput) (*ProjectEvent, error)
	DeleteEvent(ctx context.Context, id string) error
}

type ProjectService interface {
	CreateProject(ctx context.Context, input CreateProjectInput) (*Project, error)
	GetProject(ctx context.Context, id string) (*Project, error)
	ListProjects(ctx context.Context, participantUserID string) ([]Project, error)
	UpdateProject(ctx context.Context, input UpdateProjectInput) (*Project, error)
	DeleteProject(ctx context.Context, id string) error

	CreateStage(ctx context.Context, input CreateProjectStageInput) (*ProjectStage, error)
	GetStage(ctx context.Context, id string) (*ProjectStage, error)
	ListStages(ctx context.Context, projectID string) ([]ProjectStage, error)
	UpdateStage(ctx context.Context, input UpdateProjectStageInput) (*ProjectStage, error)
	DeleteStage(ctx context.Context, id string) error

	CreateMember(ctx context.Context, input CreateProjectMemberInput) (*ProjectMember, error)
	GetMember(ctx context.Context, id string) (*ProjectMember, error)
	ListMembers(ctx context.Context, projectID string) ([]ProjectMember, error)
	UpdateMember(ctx context.Context, input UpdateProjectMemberInput) (*ProjectMember, error)
	DeleteMember(ctx context.Context, id string) error

	CreateEvent(ctx context.Context, input CreateProjectEventInput) (*ProjectEvent, error)
	GetEvent(ctx context.Context, id string) (*ProjectEvent, error)
	ListEvents(ctx context.Context, projectID string) ([]ProjectEvent, error)
	UpdateEvent(ctx context.Context, input UpdateProjectEventInput) (*ProjectEvent, error)
	DeleteEvent(ctx context.Context, id string) error
}
