package grpc

import (
	"context"
	"errors"
	"time"

	projectpb "github.com/Hubcher/project-management/contracts/gen/go/project"
	"github.com/Hubcher/project-management/project-service/internal/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	projectpb.UnimplementedProjectServiceServer
	service core.ProjectService
}

func NewServer(service core.ProjectService) *Server {
	return &Server{service: service}
}

func (s *Server) Ping(_ context.Context, _ *projectpb.Empty) (*projectpb.Empty, error) {
	return &projectpb.Empty{}, nil
}

func (s *Server) CreateProject(ctx context.Context, req *projectpb.CreateProjectRequest) (*projectpb.Project, error) {
	plannedStartDate, err := parseDate(req.PlannedStartDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid planned_start_date")
	}
	plannedDeadline, err := parseDate(req.PlannedDeadline)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid planned_deadline")
	}
	actualStartDate, err := parseDate(req.ActualStartDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid actual_start_date")
	}
	actualDeadline, err := parseDate(req.ActualDeadline)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid actual_deadline")
	}

	project, err := s.service.CreateProject(ctx, core.CreateProjectInput{
		ProjectCode:      req.ProjectCode,
		Name:             req.Name,
		Description:      req.Description,
		ContractNumber:   req.ContractNumber,
		Status:           core.ProjectStatus(req.Status),
		CustomerName:     req.CustomerName,
		ManagerID:        req.ManagerId,
		PlannedStartDate: plannedStartDate,
		PlannedDeadline:  plannedDeadline,
		ActualStartDate:  actualStartDate,
		ActualDeadline:   actualDeadline,
		PlannedBudget:    req.PlannedBudget,
	})
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoProject(project), nil
}

func (s *Server) GetProject(ctx context.Context, req *projectpb.GetProjectRequest) (*projectpb.Project, error) {
	project, err := s.service.GetProject(ctx, req.Id)
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoProject(project), nil
}

func (s *Server) ListProjects(ctx context.Context, req *projectpb.ListProjectsRequest) (*projectpb.ListProjectsResponse, error) {
	projects, err := s.service.ListProjects(ctx, req.ParticipantUserId)
	if err != nil {
		return nil, mapCoreError(err)
	}
	resp := &projectpb.ListProjectsResponse{Projects: make([]*projectpb.Project, 0, len(projects))}
	for i := range projects {
		resp.Projects = append(resp.Projects, toProtoProject(&projects[i]))
	}
	return resp, nil
}

func (s *Server) UpdateProject(ctx context.Context, req *projectpb.UpdateProjectRequest) (*projectpb.Project, error) {
	plannedStartDate, err := parseDate(req.PlannedStartDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid planned_start_date")
	}
	plannedDeadline, err := parseDate(req.PlannedDeadline)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid planned_deadline")
	}
	actualStartDate, err := parseDate(req.ActualStartDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid actual_start_date")
	}
	actualDeadline, err := parseDate(req.ActualDeadline)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid actual_deadline")
	}

	project, err := s.service.UpdateProject(ctx, core.UpdateProjectInput{
		ID:               req.Id,
		ProjectCode:      req.ProjectCode,
		Name:             req.Name,
		Description:      req.Description,
		ContractNumber:   req.ContractNumber,
		Status:           core.ProjectStatus(req.Status),
		CustomerName:     req.CustomerName,
		ManagerID:        req.ManagerId,
		PlannedStartDate: plannedStartDate,
		PlannedDeadline:  plannedDeadline,
		ActualStartDate:  actualStartDate,
		ActualDeadline:   actualDeadline,
		PlannedBudget:    req.PlannedBudget,
	})
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoProject(project), nil
}

func (s *Server) DeleteProject(ctx context.Context, req *projectpb.DeleteProjectRequest) (*projectpb.Empty, error) {
	if err := s.service.DeleteProject(ctx, req.Id); err != nil {
		return nil, mapCoreError(err)
	}
	return &projectpb.Empty{}, nil
}

func (s *Server) CreateStage(ctx context.Context, req *projectpb.CreateProjectStageRequest) (*projectpb.ProjectStage, error) {
	plannedStartDate, err := parseDate(req.PlannedStartDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid planned_start_date")
	}
	plannedEndDate, err := parseDate(req.PlannedEndDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid planned_end_date")
	}
	actualStartDate, err := parseDate(req.ActualStartDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid actual_start_date")
	}
	actualEndDate, err := parseDate(req.ActualEndDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid actual_end_date")
	}

	stage, err := s.service.CreateStage(ctx, core.CreateProjectStageInput{
		ProjectID:        req.ProjectId,
		Name:             req.Name,
		Description:      req.Description,
		SequenceNumber:   req.SequenceNumber,
		Status:           core.StageStatus(req.Status),
		PlannedStartDate: plannedStartDate,
		PlannedEndDate:   plannedEndDate,
		ActualStartDate:  actualStartDate,
		ActualEndDate:    actualEndDate,
		PlannedIncome:    req.PlannedIncome,
		PlannedExpense:   req.PlannedExpense,
	})
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoStage(stage), nil
}

func (s *Server) GetStage(ctx context.Context, req *projectpb.GetProjectStageRequest) (*projectpb.ProjectStage, error) {
	stage, err := s.service.GetStage(ctx, req.Id)
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoStage(stage), nil
}

func (s *Server) ListStages(ctx context.Context, req *projectpb.ListProjectStagesRequest) (*projectpb.ListProjectStagesResponse, error) {
	stages, err := s.service.ListStages(ctx, req.ProjectId)
	if err != nil {
		return nil, mapCoreError(err)
	}
	resp := &projectpb.ListProjectStagesResponse{Stages: make([]*projectpb.ProjectStage, 0, len(stages))}
	for i := range stages {
		resp.Stages = append(resp.Stages, toProtoStage(&stages[i]))
	}
	return resp, nil
}

func (s *Server) UpdateStage(ctx context.Context, req *projectpb.UpdateProjectStageRequest) (*projectpb.ProjectStage, error) {
	plannedStartDate, err := parseDate(req.PlannedStartDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid planned_start_date")
	}
	plannedEndDate, err := parseDate(req.PlannedEndDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid planned_end_date")
	}
	actualStartDate, err := parseDate(req.ActualStartDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid actual_start_date")
	}
	actualEndDate, err := parseDate(req.ActualEndDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid actual_end_date")
	}

	stage, err := s.service.UpdateStage(ctx, core.UpdateProjectStageInput{
		ID:               req.Id,
		Name:             req.Name,
		Description:      req.Description,
		SequenceNumber:   req.SequenceNumber,
		Status:           core.StageStatus(req.Status),
		PlannedStartDate: plannedStartDate,
		PlannedEndDate:   plannedEndDate,
		ActualStartDate:  actualStartDate,
		ActualEndDate:    actualEndDate,
		PlannedIncome:    req.PlannedIncome,
		PlannedExpense:   req.PlannedExpense,
	})
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoStage(stage), nil
}

func (s *Server) DeleteStage(ctx context.Context, req *projectpb.DeleteProjectStageRequest) (*projectpb.Empty, error) {
	if err := s.service.DeleteStage(ctx, req.Id); err != nil {
		return nil, mapCoreError(err)
	}
	return &projectpb.Empty{}, nil
}

func (s *Server) CreateMember(ctx context.Context, req *projectpb.CreateProjectMemberRequest) (*projectpb.ProjectMember, error) {
	member, err := s.service.CreateMember(ctx, core.CreateProjectMemberInput{
		ProjectID:     req.ProjectId,
		UserID:        req.UserId,
		RoleInProject: req.RoleInProject,
	})
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoMember(member), nil
}

func (s *Server) GetMember(ctx context.Context, req *projectpb.GetProjectMemberRequest) (*projectpb.ProjectMember, error) {
	member, err := s.service.GetMember(ctx, req.Id)
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoMember(member), nil
}

func (s *Server) ListMembers(ctx context.Context, req *projectpb.ListProjectMembersRequest) (*projectpb.ListProjectMembersResponse, error) {
	members, err := s.service.ListMembers(ctx, req.ProjectId)
	if err != nil {
		return nil, mapCoreError(err)
	}
	resp := &projectpb.ListProjectMembersResponse{Members: make([]*projectpb.ProjectMember, 0, len(members))}
	for i := range members {
		resp.Members = append(resp.Members, toProtoMember(&members[i]))
	}
	return resp, nil
}

func (s *Server) UpdateMember(ctx context.Context, req *projectpb.UpdateProjectMemberRequest) (*projectpb.ProjectMember, error) {
	leftAt, err := parseTimestamp(req.LeftAt)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid left_at")
	}
	member, err := s.service.UpdateMember(ctx, core.UpdateProjectMemberInput{
		ID:            req.Id,
		RoleInProject: req.RoleInProject,
		IsActive:      req.IsActive,
		LeftAt:        leftAt,
	})
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoMember(member), nil
}

func (s *Server) DeleteMember(ctx context.Context, req *projectpb.DeleteProjectMemberRequest) (*projectpb.Empty, error) {
	if err := s.service.DeleteMember(ctx, req.Id); err != nil {
		return nil, mapCoreError(err)
	}
	return &projectpb.Empty{}, nil
}

func (s *Server) CreateEvent(ctx context.Context, req *projectpb.CreateProjectEventRequest) (*projectpb.ProjectEvent, error) {
	plannedDate, err := parseDate(req.PlannedDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid planned_date")
	}
	actualDate, err := parseDate(req.ActualDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid actual_date")
	}
	projectEvent, err := s.service.CreateEvent(ctx, core.CreateProjectEventInput{
		ProjectID:   req.ProjectId,
		StageID:     req.StageId,
		Name:        req.Name,
		Description: req.Description,
		PlannedDate: plannedDate,
		ActualDate:  actualDate,
		Status:      core.EventStatus(req.Status),
	})
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoEvent(projectEvent), nil
}

func (s *Server) GetEvent(ctx context.Context, req *projectpb.GetProjectEventRequest) (*projectpb.ProjectEvent, error) {
	projectEvent, err := s.service.GetEvent(ctx, req.Id)
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoEvent(projectEvent), nil
}

func (s *Server) ListEvents(ctx context.Context, req *projectpb.ListProjectEventsRequest) (*projectpb.ListProjectEventsResponse, error) {
	events, err := s.service.ListEvents(ctx, req.ProjectId)
	if err != nil {
		return nil, mapCoreError(err)
	}
	resp := &projectpb.ListProjectEventsResponse{Events: make([]*projectpb.ProjectEvent, 0, len(events))}
	for i := range events {
		resp.Events = append(resp.Events, toProtoEvent(&events[i]))
	}
	return resp, nil
}

func (s *Server) UpdateEvent(ctx context.Context, req *projectpb.UpdateProjectEventRequest) (*projectpb.ProjectEvent, error) {
	plannedDate, err := parseDate(req.PlannedDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid planned_date")
	}
	actualDate, err := parseDate(req.ActualDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid actual_date")
	}
	projectEvent, err := s.service.UpdateEvent(ctx, core.UpdateProjectEventInput{
		ID:          req.Id,
		StageID:     req.StageId,
		Name:        req.Name,
		Description: req.Description,
		PlannedDate: plannedDate,
		ActualDate:  actualDate,
		Status:      core.EventStatus(req.Status),
	})
	if err != nil {
		return nil, mapCoreError(err)
	}
	return toProtoEvent(projectEvent), nil
}

func (s *Server) DeleteEvent(ctx context.Context, req *projectpb.DeleteProjectEventRequest) (*projectpb.Empty, error) {
	if err := s.service.DeleteEvent(ctx, req.Id); err != nil {
		return nil, mapCoreError(err)
	}
	return &projectpb.Empty{}, nil
}

func parseDate(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, err
	}
	day := parsed.UTC()
	normalized := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	return &normalized, nil
}

func parseTimestamp(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	if ts, err := time.Parse(time.RFC3339, value); err == nil {
		normalized := ts.UTC()
		return &normalized, nil
	}
	return parseDate(value)
}

func toProtoProject(project *core.Project) *projectpb.Project {
	return &projectpb.Project{
		ID:               project.ID,
		ProjectCode:      project.ProjectCode,
		Name:             project.Name,
		Description:      project.Description,
		ContractNumber:   project.ContractNumber,
		Status:           string(project.Status),
		CustomerName:     project.CustomerName,
		ManagerId:        project.ManagerID,
		PlannedStartDate: formatDate(project.PlannedStartDate),
		PlannedDeadline:  formatDate(project.PlannedDeadline),
		ActualStartDate:  formatDate(project.ActualStartDate),
		ActualDeadline:   formatDate(project.ActualDeadline),
		PlannedBudget:    project.PlannedBudget,
		CreatedAt:        formatTimestamp(&project.CreatedAt),
		UpdatedAt:        formatTimestamp(&project.UpdatedAt),
	}
}

func toProtoStage(stage *core.ProjectStage) *projectpb.ProjectStage {
	return &projectpb.ProjectStage{
		ID:               stage.ID,
		ProjectId:        stage.ProjectID,
		Name:             stage.Name,
		Description:      stage.Description,
		SequenceNumber:   stage.SequenceNumber,
		Status:           string(stage.Status),
		PlannedStartDate: formatDate(stage.PlannedStartDate),
		PlannedEndDate:   formatDate(stage.PlannedEndDate),
		ActualStartDate:  formatDate(stage.ActualStartDate),
		ActualEndDate:    formatDate(stage.ActualEndDate),
		PlannedIncome:    stage.PlannedIncome,
		PlannedExpense:   stage.PlannedExpense,
		CreatedAt:        formatTimestamp(&stage.CreatedAt),
		UpdatedAt:        formatTimestamp(&stage.UpdatedAt),
	}
}

func toProtoMember(member *core.ProjectMember) *projectpb.ProjectMember {
	return &projectpb.ProjectMember{
		ID:            member.ID,
		ProjectId:     member.ProjectID,
		UserId:        member.UserID,
		RoleInProject: member.RoleInProject,
		IsActive:      member.IsActive,
		JoinedAt:      formatTimestamp(&member.JoinedAt),
		LeftAt:        formatTimestamp(member.LeftAt),
	}
}

func toProtoEvent(projectEvent *core.ProjectEvent) *projectpb.ProjectEvent {
	return &projectpb.ProjectEvent{
		ID:          projectEvent.ID,
		ProjectId:   projectEvent.ProjectID,
		StageId:     projectEvent.StageID,
		Name:        projectEvent.Name,
		Description: projectEvent.Description,
		PlannedDate: formatDate(projectEvent.PlannedDate),
		ActualDate:  formatDate(projectEvent.ActualDate),
		Status:      string(projectEvent.Status),
		CreatedAt:   formatTimestamp(&projectEvent.CreatedAt),
		UpdatedAt:   formatTimestamp(&projectEvent.UpdatedAt),
	}
}

func formatDate(value *time.Time) string {
	if value == nil || value.IsZero() {
		return ""
	}
	return value.UTC().Format("2006-01-02")
}

func formatTimestamp(value *time.Time) string {
	if value == nil || value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339)
}

func mapCoreError(err error) error {
	switch {
	case errors.Is(err, core.ErrInvalidProject), errors.Is(err, core.ErrInvalidStage), errors.Is(err, core.ErrInvalidMember), errors.Is(err, core.ErrInvalidEvent):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, core.ErrProjectNotFound), errors.Is(err, core.ErrStageNotFound), errors.Is(err, core.ErrMemberNotFound), errors.Is(err, core.ErrEventNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, core.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
