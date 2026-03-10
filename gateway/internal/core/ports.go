package core

import (
	"context"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

//type ProjectService interface {
//	CreateProject(ctx context.Context, req *projectpb.CreateProjectRequest) (*projectpb.Project, error)
//	GetProject(ctx context.Context, req *projectpb.GetProjectRequest) (*projectpb.Project, error)
//	ListProject(ctx context.Context, req *projectpb.ListProjectsRequest) (*projectpb.ListProjectsResponse, error)
//	UpdateProject(ctx context.Context, req *projectpb.UpdateProjectRequest) (*projectpb.Project, error)
//	DeleteProject(ctx context.Context, req *projectpb.DeleteProjectRequest) (*empty.Empty, error)
//}
