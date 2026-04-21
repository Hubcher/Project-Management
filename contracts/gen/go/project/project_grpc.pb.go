package projectpb

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ProjectService_ServiceName                  = "project.ProjectService"
	ProjectService_Ping_FullMethodName         = "/project.ProjectService/Ping"
	ProjectService_CreateProject_FullMethodName = "/project.ProjectService/CreateProject"
	ProjectService_GetProject_FullMethodName    = "/project.ProjectService/GetProject"
	ProjectService_ListProjects_FullMethodName  = "/project.ProjectService/ListProjects"
	ProjectService_UpdateProject_FullMethodName = "/project.ProjectService/UpdateProject"
	ProjectService_DeleteProject_FullMethodName = "/project.ProjectService/DeleteProject"
	ProjectService_CreateStage_FullMethodName   = "/project.ProjectService/CreateStage"
	ProjectService_GetStage_FullMethodName      = "/project.ProjectService/GetStage"
	ProjectService_ListStages_FullMethodName    = "/project.ProjectService/ListStages"
	ProjectService_UpdateStage_FullMethodName   = "/project.ProjectService/UpdateStage"
	ProjectService_DeleteStage_FullMethodName   = "/project.ProjectService/DeleteStage"
	ProjectService_CreateMember_FullMethodName  = "/project.ProjectService/CreateMember"
	ProjectService_GetMember_FullMethodName     = "/project.ProjectService/GetMember"
	ProjectService_ListMembers_FullMethodName   = "/project.ProjectService/ListMembers"
	ProjectService_UpdateMember_FullMethodName  = "/project.ProjectService/UpdateMember"
	ProjectService_DeleteMember_FullMethodName  = "/project.ProjectService/DeleteMember"
	ProjectService_CreateEvent_FullMethodName   = "/project.ProjectService/CreateEvent"
	ProjectService_GetEvent_FullMethodName      = "/project.ProjectService/GetEvent"
	ProjectService_ListEvents_FullMethodName    = "/project.ProjectService/ListEvents"
	ProjectService_UpdateEvent_FullMethodName   = "/project.ProjectService/UpdateEvent"
	ProjectService_DeleteEvent_FullMethodName   = "/project.ProjectService/DeleteEvent"
)

type ProjectServiceClient interface {
	Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	CreateProject(ctx context.Context, in *CreateProjectRequest, opts ...grpc.CallOption) (*Project, error)
	GetProject(ctx context.Context, in *GetProjectRequest, opts ...grpc.CallOption) (*Project, error)
	ListProjects(ctx context.Context, in *ListProjectsRequest, opts ...grpc.CallOption) (*ListProjectsResponse, error)
	UpdateProject(ctx context.Context, in *UpdateProjectRequest, opts ...grpc.CallOption) (*Project, error)
	DeleteProject(ctx context.Context, in *DeleteProjectRequest, opts ...grpc.CallOption) (*Empty, error)
	CreateStage(ctx context.Context, in *CreateProjectStageRequest, opts ...grpc.CallOption) (*ProjectStage, error)
	GetStage(ctx context.Context, in *GetProjectStageRequest, opts ...grpc.CallOption) (*ProjectStage, error)
	ListStages(ctx context.Context, in *ListProjectStagesRequest, opts ...grpc.CallOption) (*ListProjectStagesResponse, error)
	UpdateStage(ctx context.Context, in *UpdateProjectStageRequest, opts ...grpc.CallOption) (*ProjectStage, error)
	DeleteStage(ctx context.Context, in *DeleteProjectStageRequest, opts ...grpc.CallOption) (*Empty, error)
	CreateMember(ctx context.Context, in *CreateProjectMemberRequest, opts ...grpc.CallOption) (*ProjectMember, error)
	GetMember(ctx context.Context, in *GetProjectMemberRequest, opts ...grpc.CallOption) (*ProjectMember, error)
	ListMembers(ctx context.Context, in *ListProjectMembersRequest, opts ...grpc.CallOption) (*ListProjectMembersResponse, error)
	UpdateMember(ctx context.Context, in *UpdateProjectMemberRequest, opts ...grpc.CallOption) (*ProjectMember, error)
	DeleteMember(ctx context.Context, in *DeleteProjectMemberRequest, opts ...grpc.CallOption) (*Empty, error)
	CreateEvent(ctx context.Context, in *CreateProjectEventRequest, opts ...grpc.CallOption) (*ProjectEvent, error)
	GetEvent(ctx context.Context, in *GetProjectEventRequest, opts ...grpc.CallOption) (*ProjectEvent, error)
	ListEvents(ctx context.Context, in *ListProjectEventsRequest, opts ...grpc.CallOption) (*ListProjectEventsResponse, error)
	UpdateEvent(ctx context.Context, in *UpdateProjectEventRequest, opts ...grpc.CallOption) (*ProjectEvent, error)
	DeleteEvent(ctx context.Context, in *DeleteProjectEventRequest, opts ...grpc.CallOption) (*Empty, error)
}

type projectServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProjectServiceClient(cc grpc.ClientConnInterface) ProjectServiceClient {
	return &projectServiceClient{cc: cc}
}

func invoke[Resp any](ctx context.Context, cc grpc.ClientConnInterface, method string, in any, out Resp, opts ...grpc.CallOption) (Resp, error) {
	if err := cc.Invoke(ctx, method, in, out, opts...); err != nil {
		var zero Resp
		return zero, err
	}
	return out, nil
}

func (c *projectServiceClient) Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	return invoke(ctx, c.cc, ProjectService_Ping_FullMethodName, in, new(Empty), opts...)
}
func (c *projectServiceClient) CreateProject(ctx context.Context, in *CreateProjectRequest, opts ...grpc.CallOption) (*Project, error) {
	return invoke(ctx, c.cc, ProjectService_CreateProject_FullMethodName, in, new(Project), opts...)
}
func (c *projectServiceClient) GetProject(ctx context.Context, in *GetProjectRequest, opts ...grpc.CallOption) (*Project, error) {
	return invoke(ctx, c.cc, ProjectService_GetProject_FullMethodName, in, new(Project), opts...)
}
func (c *projectServiceClient) ListProjects(ctx context.Context, in *ListProjectsRequest, opts ...grpc.CallOption) (*ListProjectsResponse, error) {
	return invoke(ctx, c.cc, ProjectService_ListProjects_FullMethodName, in, new(ListProjectsResponse), opts...)
}
func (c *projectServiceClient) UpdateProject(ctx context.Context, in *UpdateProjectRequest, opts ...grpc.CallOption) (*Project, error) {
	return invoke(ctx, c.cc, ProjectService_UpdateProject_FullMethodName, in, new(Project), opts...)
}
func (c *projectServiceClient) DeleteProject(ctx context.Context, in *DeleteProjectRequest, opts ...grpc.CallOption) (*Empty, error) {
	return invoke(ctx, c.cc, ProjectService_DeleteProject_FullMethodName, in, new(Empty), opts...)
}
func (c *projectServiceClient) CreateStage(ctx context.Context, in *CreateProjectStageRequest, opts ...grpc.CallOption) (*ProjectStage, error) {
	return invoke(ctx, c.cc, ProjectService_CreateStage_FullMethodName, in, new(ProjectStage), opts...)
}
func (c *projectServiceClient) GetStage(ctx context.Context, in *GetProjectStageRequest, opts ...grpc.CallOption) (*ProjectStage, error) {
	return invoke(ctx, c.cc, ProjectService_GetStage_FullMethodName, in, new(ProjectStage), opts...)
}
func (c *projectServiceClient) ListStages(ctx context.Context, in *ListProjectStagesRequest, opts ...grpc.CallOption) (*ListProjectStagesResponse, error) {
	return invoke(ctx, c.cc, ProjectService_ListStages_FullMethodName, in, new(ListProjectStagesResponse), opts...)
}
func (c *projectServiceClient) UpdateStage(ctx context.Context, in *UpdateProjectStageRequest, opts ...grpc.CallOption) (*ProjectStage, error) {
	return invoke(ctx, c.cc, ProjectService_UpdateStage_FullMethodName, in, new(ProjectStage), opts...)
}
func (c *projectServiceClient) DeleteStage(ctx context.Context, in *DeleteProjectStageRequest, opts ...grpc.CallOption) (*Empty, error) {
	return invoke(ctx, c.cc, ProjectService_DeleteStage_FullMethodName, in, new(Empty), opts...)
}
func (c *projectServiceClient) CreateMember(ctx context.Context, in *CreateProjectMemberRequest, opts ...grpc.CallOption) (*ProjectMember, error) {
	return invoke(ctx, c.cc, ProjectService_CreateMember_FullMethodName, in, new(ProjectMember), opts...)
}
func (c *projectServiceClient) GetMember(ctx context.Context, in *GetProjectMemberRequest, opts ...grpc.CallOption) (*ProjectMember, error) {
	return invoke(ctx, c.cc, ProjectService_GetMember_FullMethodName, in, new(ProjectMember), opts...)
}
func (c *projectServiceClient) ListMembers(ctx context.Context, in *ListProjectMembersRequest, opts ...grpc.CallOption) (*ListProjectMembersResponse, error) {
	return invoke(ctx, c.cc, ProjectService_ListMembers_FullMethodName, in, new(ListProjectMembersResponse), opts...)
}
func (c *projectServiceClient) UpdateMember(ctx context.Context, in *UpdateProjectMemberRequest, opts ...grpc.CallOption) (*ProjectMember, error) {
	return invoke(ctx, c.cc, ProjectService_UpdateMember_FullMethodName, in, new(ProjectMember), opts...)
}
func (c *projectServiceClient) DeleteMember(ctx context.Context, in *DeleteProjectMemberRequest, opts ...grpc.CallOption) (*Empty, error) {
	return invoke(ctx, c.cc, ProjectService_DeleteMember_FullMethodName, in, new(Empty), opts...)
}
func (c *projectServiceClient) CreateEvent(ctx context.Context, in *CreateProjectEventRequest, opts ...grpc.CallOption) (*ProjectEvent, error) {
	return invoke(ctx, c.cc, ProjectService_CreateEvent_FullMethodName, in, new(ProjectEvent), opts...)
}
func (c *projectServiceClient) GetEvent(ctx context.Context, in *GetProjectEventRequest, opts ...grpc.CallOption) (*ProjectEvent, error) {
	return invoke(ctx, c.cc, ProjectService_GetEvent_FullMethodName, in, new(ProjectEvent), opts...)
}
func (c *projectServiceClient) ListEvents(ctx context.Context, in *ListProjectEventsRequest, opts ...grpc.CallOption) (*ListProjectEventsResponse, error) {
	return invoke(ctx, c.cc, ProjectService_ListEvents_FullMethodName, in, new(ListProjectEventsResponse), opts...)
}
func (c *projectServiceClient) UpdateEvent(ctx context.Context, in *UpdateProjectEventRequest, opts ...grpc.CallOption) (*ProjectEvent, error) {
	return invoke(ctx, c.cc, ProjectService_UpdateEvent_FullMethodName, in, new(ProjectEvent), opts...)
}
func (c *projectServiceClient) DeleteEvent(ctx context.Context, in *DeleteProjectEventRequest, opts ...grpc.CallOption) (*Empty, error) {
	return invoke(ctx, c.cc, ProjectService_DeleteEvent_FullMethodName, in, new(Empty), opts...)
}

type ProjectServiceServer interface {
	Ping(context.Context, *Empty) (*Empty, error)
	CreateProject(context.Context, *CreateProjectRequest) (*Project, error)
	GetProject(context.Context, *GetProjectRequest) (*Project, error)
	ListProjects(context.Context, *ListProjectsRequest) (*ListProjectsResponse, error)
	UpdateProject(context.Context, *UpdateProjectRequest) (*Project, error)
	DeleteProject(context.Context, *DeleteProjectRequest) (*Empty, error)
	CreateStage(context.Context, *CreateProjectStageRequest) (*ProjectStage, error)
	GetStage(context.Context, *GetProjectStageRequest) (*ProjectStage, error)
	ListStages(context.Context, *ListProjectStagesRequest) (*ListProjectStagesResponse, error)
	UpdateStage(context.Context, *UpdateProjectStageRequest) (*ProjectStage, error)
	DeleteStage(context.Context, *DeleteProjectStageRequest) (*Empty, error)
	CreateMember(context.Context, *CreateProjectMemberRequest) (*ProjectMember, error)
	GetMember(context.Context, *GetProjectMemberRequest) (*ProjectMember, error)
	ListMembers(context.Context, *ListProjectMembersRequest) (*ListProjectMembersResponse, error)
	UpdateMember(context.Context, *UpdateProjectMemberRequest) (*ProjectMember, error)
	DeleteMember(context.Context, *DeleteProjectMemberRequest) (*Empty, error)
	CreateEvent(context.Context, *CreateProjectEventRequest) (*ProjectEvent, error)
	GetEvent(context.Context, *GetProjectEventRequest) (*ProjectEvent, error)
	ListEvents(context.Context, *ListProjectEventsRequest) (*ListProjectEventsResponse, error)
	UpdateEvent(context.Context, *UpdateProjectEventRequest) (*ProjectEvent, error)
	DeleteEvent(context.Context, *DeleteProjectEventRequest) (*Empty, error)
}

type UnimplementedProjectServiceServer struct{}

func (UnimplementedProjectServiceServer) Ping(context.Context, *Empty) (*Empty, error) { return nil, status.Error(codes.Unimplemented, "method Ping not implemented") }
func (UnimplementedProjectServiceServer) CreateProject(context.Context, *CreateProjectRequest) (*Project, error) {
	return nil, status.Error(codes.Unimplemented, "method CreateProject not implemented")
}
func (UnimplementedProjectServiceServer) GetProject(context.Context, *GetProjectRequest) (*Project, error) {
	return nil, status.Error(codes.Unimplemented, "method GetProject not implemented")
}
func (UnimplementedProjectServiceServer) ListProjects(context.Context, *ListProjectsRequest) (*ListProjectsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method ListProjects not implemented")
}
func (UnimplementedProjectServiceServer) UpdateProject(context.Context, *UpdateProjectRequest) (*Project, error) {
	return nil, status.Error(codes.Unimplemented, "method UpdateProject not implemented")
}
func (UnimplementedProjectServiceServer) DeleteProject(context.Context, *DeleteProjectRequest) (*Empty, error) {
	return nil, status.Error(codes.Unimplemented, "method DeleteProject not implemented")
}
func (UnimplementedProjectServiceServer) CreateStage(context.Context, *CreateProjectStageRequest) (*ProjectStage, error) {
	return nil, status.Error(codes.Unimplemented, "method CreateStage not implemented")
}
func (UnimplementedProjectServiceServer) GetStage(context.Context, *GetProjectStageRequest) (*ProjectStage, error) {
	return nil, status.Error(codes.Unimplemented, "method GetStage not implemented")
}
func (UnimplementedProjectServiceServer) ListStages(context.Context, *ListProjectStagesRequest) (*ListProjectStagesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method ListStages not implemented")
}
func (UnimplementedProjectServiceServer) UpdateStage(context.Context, *UpdateProjectStageRequest) (*ProjectStage, error) {
	return nil, status.Error(codes.Unimplemented, "method UpdateStage not implemented")
}
func (UnimplementedProjectServiceServer) DeleteStage(context.Context, *DeleteProjectStageRequest) (*Empty, error) {
	return nil, status.Error(codes.Unimplemented, "method DeleteStage not implemented")
}
func (UnimplementedProjectServiceServer) CreateMember(context.Context, *CreateProjectMemberRequest) (*ProjectMember, error) {
	return nil, status.Error(codes.Unimplemented, "method CreateMember not implemented")
}
func (UnimplementedProjectServiceServer) GetMember(context.Context, *GetProjectMemberRequest) (*ProjectMember, error) {
	return nil, status.Error(codes.Unimplemented, "method GetMember not implemented")
}
func (UnimplementedProjectServiceServer) ListMembers(context.Context, *ListProjectMembersRequest) (*ListProjectMembersResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method ListMembers not implemented")
}
func (UnimplementedProjectServiceServer) UpdateMember(context.Context, *UpdateProjectMemberRequest) (*ProjectMember, error) {
	return nil, status.Error(codes.Unimplemented, "method UpdateMember not implemented")
}
func (UnimplementedProjectServiceServer) DeleteMember(context.Context, *DeleteProjectMemberRequest) (*Empty, error) {
	return nil, status.Error(codes.Unimplemented, "method DeleteMember not implemented")
}
func (UnimplementedProjectServiceServer) CreateEvent(context.Context, *CreateProjectEventRequest) (*ProjectEvent, error) {
	return nil, status.Error(codes.Unimplemented, "method CreateEvent not implemented")
}
func (UnimplementedProjectServiceServer) GetEvent(context.Context, *GetProjectEventRequest) (*ProjectEvent, error) {
	return nil, status.Error(codes.Unimplemented, "method GetEvent not implemented")
}
func (UnimplementedProjectServiceServer) ListEvents(context.Context, *ListProjectEventsRequest) (*ListProjectEventsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method ListEvents not implemented")
}
func (UnimplementedProjectServiceServer) UpdateEvent(context.Context, *UpdateProjectEventRequest) (*ProjectEvent, error) {
	return nil, status.Error(codes.Unimplemented, "method UpdateEvent not implemented")
}
func (UnimplementedProjectServiceServer) DeleteEvent(context.Context, *DeleteProjectEventRequest) (*Empty, error) {
	return nil, status.Error(codes.Unimplemented, "method DeleteEvent not implemented")
}

func RegisterProjectServiceServer(s grpc.ServiceRegistrar, srv ProjectServiceServer) {
	s.RegisterService(&ProjectService_ServiceDesc, srv)
}

func handleUnary(ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor, fullMethod string, in any, srv interface{}, call func(context.Context, any) (any, error)) (any, error) {
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return call(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: fullMethod}
	handler := func(ctx context.Context, req any) (any, error) {
		return call(ctx, req)
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_Ping_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ProjectService_Ping_FullMethodName, new(Empty), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ProjectServiceServer).Ping(ctx, req.(*Empty))
	})
}
func _ProjectService_CreateProject_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ProjectService_CreateProject_FullMethodName, new(CreateProjectRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ProjectServiceServer).CreateProject(ctx, req.(*CreateProjectRequest))
	})
}
func _ProjectService_GetProject_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ProjectService_GetProject_FullMethodName, new(GetProjectRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ProjectServiceServer).GetProject(ctx, req.(*GetProjectRequest))
	})
}
func _ProjectService_ListProjects_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ProjectService_ListProjects_FullMethodName, new(ListProjectsRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ProjectServiceServer).ListProjects(ctx, req.(*ListProjectsRequest))
	})
}
func _ProjectService_UpdateProject_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ProjectService_UpdateProject_FullMethodName, new(UpdateProjectRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ProjectServiceServer).UpdateProject(ctx, req.(*UpdateProjectRequest))
	})
}
func _ProjectService_DeleteProject_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ProjectService_DeleteProject_FullMethodName, new(DeleteProjectRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ProjectServiceServer).DeleteProject(ctx, req.(*DeleteProjectRequest))
	})
}
func _ProjectService_CreateStage_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ProjectService_CreateStage_FullMethodName, new(CreateProjectStageRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ProjectServiceServer).CreateStage(ctx, req.(*CreateProjectStageRequest))
	})
}
func _ProjectService_GetStage_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ProjectService_GetStage_FullMethodName, new(GetProjectStageRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ProjectServiceServer).GetStage(ctx, req.(*GetProjectStageRequest))
	})
}
func _ProjectService_ListStages_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ProjectService_ListStages_FullMethodName, new(ListProjectStagesRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ProjectServiceServer).ListStages(ctx, req.(*ListProjectStagesRequest))
	})
}
func _ProjectService_UpdateStage_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ProjectService_UpdateStage_FullMethodName, new(UpdateProjectStageRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ProjectServiceServer).UpdateStage(ctx, req.(*UpdateProjectStageRequest))
	})
}
func _ProjectService_DeleteStage_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ProjectService_DeleteStage_FullMethodName, new(DeleteProjectStageRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ProjectServiceServer).DeleteStage(ctx, req.(*DeleteProjectStageRequest))
	})
}
func _ProjectService_CreateMember_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ProjectService_CreateMember_FullMethodName, new(CreateProjectMemberRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ProjectServiceServer).CreateMember(ctx, req.(*CreateProjectMemberRequest))
	})
}
func _ProjectService_GetMember_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ProjectService_GetMember_FullMethodName, new(GetProjectMemberRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ProjectServiceServer).GetMember(ctx, req.(*GetProjectMemberRequest))
	})
}
func _ProjectService_ListMembers_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ProjectService_ListMembers_FullMethodName, new(ListProjectMembersRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ProjectServiceServer).ListMembers(ctx, req.(*ListProjectMembersRequest))
	})
}
func _ProjectService_UpdateMember_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ProjectService_UpdateMember_FullMethodName, new(UpdateProjectMemberRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ProjectServiceServer).UpdateMember(ctx, req.(*UpdateProjectMemberRequest))
	})
}
func _ProjectService_DeleteMember_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ProjectService_DeleteMember_FullMethodName, new(DeleteProjectMemberRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ProjectServiceServer).DeleteMember(ctx, req.(*DeleteProjectMemberRequest))
	})
}
func _ProjectService_CreateEvent_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ProjectService_CreateEvent_FullMethodName, new(CreateProjectEventRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ProjectServiceServer).CreateEvent(ctx, req.(*CreateProjectEventRequest))
	})
}
func _ProjectService_GetEvent_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ProjectService_GetEvent_FullMethodName, new(GetProjectEventRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ProjectServiceServer).GetEvent(ctx, req.(*GetProjectEventRequest))
	})
}
func _ProjectService_ListEvents_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ProjectService_ListEvents_FullMethodName, new(ListProjectEventsRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ProjectServiceServer).ListEvents(ctx, req.(*ListProjectEventsRequest))
	})
}
func _ProjectService_UpdateEvent_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ProjectService_UpdateEvent_FullMethodName, new(UpdateProjectEventRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ProjectServiceServer).UpdateEvent(ctx, req.(*UpdateProjectEventRequest))
	})
}
func _ProjectService_DeleteEvent_Handler(srv interface{}, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	return handleUnary(ctx, dec, interceptor, ProjectService_DeleteEvent_FullMethodName, new(DeleteProjectEventRequest), srv, func(ctx context.Context, req any) (any, error) {
		return srv.(ProjectServiceServer).DeleteEvent(ctx, req.(*DeleteProjectEventRequest))
	})
}

var ProjectService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: ProjectService_ServiceName,
	HandlerType: (*ProjectServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "Ping", Handler: _ProjectService_Ping_Handler},
		{MethodName: "CreateProject", Handler: _ProjectService_CreateProject_Handler},
		{MethodName: "GetProject", Handler: _ProjectService_GetProject_Handler},
		{MethodName: "ListProjects", Handler: _ProjectService_ListProjects_Handler},
		{MethodName: "UpdateProject", Handler: _ProjectService_UpdateProject_Handler},
		{MethodName: "DeleteProject", Handler: _ProjectService_DeleteProject_Handler},
		{MethodName: "CreateStage", Handler: _ProjectService_CreateStage_Handler},
		{MethodName: "GetStage", Handler: _ProjectService_GetStage_Handler},
		{MethodName: "ListStages", Handler: _ProjectService_ListStages_Handler},
		{MethodName: "UpdateStage", Handler: _ProjectService_UpdateStage_Handler},
		{MethodName: "DeleteStage", Handler: _ProjectService_DeleteStage_Handler},
		{MethodName: "CreateMember", Handler: _ProjectService_CreateMember_Handler},
		{MethodName: "GetMember", Handler: _ProjectService_GetMember_Handler},
		{MethodName: "ListMembers", Handler: _ProjectService_ListMembers_Handler},
		{MethodName: "UpdateMember", Handler: _ProjectService_UpdateMember_Handler},
		{MethodName: "DeleteMember", Handler: _ProjectService_DeleteMember_Handler},
		{MethodName: "CreateEvent", Handler: _ProjectService_CreateEvent_Handler},
		{MethodName: "GetEvent", Handler: _ProjectService_GetEvent_Handler},
		{MethodName: "ListEvents", Handler: _ProjectService_ListEvents_Handler},
		{MethodName: "UpdateEvent", Handler: _ProjectService_UpdateEvent_Handler},
		{MethodName: "DeleteEvent", Handler: _ProjectService_DeleteEvent_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "project/project.proto",
}
