package project

import (
	"context"
	"errors"
	"time"

	projectpb "github.com/Hubcher/project-management/contracts/gen/go/project"
	"github.com/Hubcher/project-management/gateway/internal/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client projectpb.ProjectServiceClient
	conn   *grpc.ClientConn
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.CallContentSubtype(projectpb.JSONCodecName)),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.Config{BaseDelay: time.Second, Multiplier: 1.6, MaxDelay: 5 * time.Second},
			MinConnectTimeout: 10 * time.Second,
		}),
	)
	if err != nil {
		return nil, err
	}
	conn.Connect()
	return &Client{client: projectpb.NewProjectServiceClient(conn), conn: conn}, nil
}

func (c *Client) Close() error { return c.conn.Close() }

func (c *Client) Ping(ctx context.Context) error {
	_, err := callWithTimeout(ctx, func(callCtx context.Context) (*projectpb.Empty, error) {
		return c.client.Ping(callCtx, &projectpb.Empty{})
	})
	return err
}

func (c *Client) CreateProject(ctx context.Context, input core.CreateProjectInput) (*core.Project, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*projectpb.Project, error) {
		return c.client.CreateProject(callCtx, &projectpb.CreateProjectRequest{
			ProjectCode:      input.ProjectCode,
			Name:             input.Name,
			Description:      input.Description,
			ContractNumber:   input.ContractNumber,
			Status:           input.Status,
			CustomerName:     input.CustomerName,
			ManagerId:        input.ManagerID,
			PlannedStartDate: input.PlannedStartDate,
			PlannedDeadline:  input.PlannedDeadline,
			ActualStartDate:  input.ActualStartDate,
			ActualDeadline:   input.ActualDeadline,
			PlannedBudget:    input.PlannedBudget,
		})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoProject(resp)
}

func (c *Client) GetProject(ctx context.Context, id string) (*core.Project, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*projectpb.Project, error) {
		return c.client.GetProject(callCtx, &projectpb.GetProjectRequest{Id: id})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoProject(resp)
}

func (c *Client) ListProjects(ctx context.Context, participantUserID string) ([]core.Project, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*projectpb.ListProjectsResponse, error) {
		return c.client.ListProjects(callCtx, &projectpb.ListProjectsRequest{ParticipantUserId: participantUserID})
	})
	if err != nil {
		return nil, err
	}
	items := make([]core.Project, 0, len(resp.Projects))
	for _, project := range resp.Projects {
		item, convErr := fromProtoProject(project)
		if convErr != nil {
			return nil, convErr
		}
		items = append(items, *item)
	}
	return items, nil
}

func (c *Client) UpdateProject(ctx context.Context, input core.UpdateProjectInput) (*core.Project, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*projectpb.Project, error) {
		return c.client.UpdateProject(callCtx, &projectpb.UpdateProjectRequest{
			Id:               input.ID,
			ProjectCode:      input.ProjectCode,
			Name:             input.Name,
			Description:      input.Description,
			ContractNumber:   input.ContractNumber,
			Status:           input.Status,
			CustomerName:     input.CustomerName,
			ManagerId:        input.ManagerID,
			PlannedStartDate: input.PlannedStartDate,
			PlannedDeadline:  input.PlannedDeadline,
			ActualStartDate:  input.ActualStartDate,
			ActualDeadline:   input.ActualDeadline,
			PlannedBudget:    input.PlannedBudget,
		})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoProject(resp)
}

func (c *Client) DeleteProject(ctx context.Context, id string) error {
	_, err := callWithTimeout(ctx, func(callCtx context.Context) (*projectpb.Empty, error) {
		return c.client.DeleteProject(callCtx, &projectpb.DeleteProjectRequest{Id: id})
	})
	return err
}

func (c *Client) CreateStage(ctx context.Context, input core.CreateProjectStageInput) (*core.ProjectStage, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*projectpb.ProjectStage, error) {
		return c.client.CreateStage(callCtx, &projectpb.CreateProjectStageRequest{
			ProjectId:        input.ProjectID,
			Name:             input.Name,
			Description:      input.Description,
			SequenceNumber:   input.SequenceNumber,
			Status:           input.Status,
			PlannedStartDate: input.PlannedStartDate,
			PlannedEndDate:   input.PlannedEndDate,
			ActualStartDate:  input.ActualStartDate,
			ActualEndDate:    input.ActualEndDate,
			PlannedIncome:    input.PlannedIncome,
			PlannedExpense:   input.PlannedExpense,
		})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoStage(resp)
}

func (c *Client) GetStage(ctx context.Context, id string) (*core.ProjectStage, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*projectpb.ProjectStage, error) {
		return c.client.GetStage(callCtx, &projectpb.GetProjectStageRequest{Id: id})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoStage(resp)
}

func (c *Client) ListStages(ctx context.Context, projectID string) ([]core.ProjectStage, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*projectpb.ListProjectStagesResponse, error) {
		return c.client.ListStages(callCtx, &projectpb.ListProjectStagesRequest{ProjectId: projectID})
	})
	if err != nil {
		return nil, err
	}
	items := make([]core.ProjectStage, 0, len(resp.Stages))
	for _, stage := range resp.Stages {
		item, convErr := fromProtoStage(stage)
		if convErr != nil {
			return nil, convErr
		}
		items = append(items, *item)
	}
	return items, nil
}

func (c *Client) UpdateStage(ctx context.Context, input core.UpdateProjectStageInput) (*core.ProjectStage, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*projectpb.ProjectStage, error) {
		return c.client.UpdateStage(callCtx, &projectpb.UpdateProjectStageRequest{
			Id:               input.ID,
			Name:             input.Name,
			Description:      input.Description,
			SequenceNumber:   input.SequenceNumber,
			Status:           input.Status,
			PlannedStartDate: input.PlannedStartDate,
			PlannedEndDate:   input.PlannedEndDate,
			ActualStartDate:  input.ActualStartDate,
			ActualEndDate:    input.ActualEndDate,
			PlannedIncome:    input.PlannedIncome,
			PlannedExpense:   input.PlannedExpense,
		})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoStage(resp)
}

func (c *Client) DeleteStage(ctx context.Context, id string) error {
	_, err := callWithTimeout(ctx, func(callCtx context.Context) (*projectpb.Empty, error) {
		return c.client.DeleteStage(callCtx, &projectpb.DeleteProjectStageRequest{Id: id})
	})
	return err
}

func (c *Client) CreateMember(ctx context.Context, input core.CreateProjectMemberInput) (*core.ProjectMember, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*projectpb.ProjectMember, error) {
		return c.client.CreateMember(callCtx, &projectpb.CreateProjectMemberRequest{
			ProjectId:     input.ProjectID,
			UserId:        input.UserID,
			RoleInProject: input.RoleInProject,
		})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoMember(resp)
}

func (c *Client) GetMember(ctx context.Context, id string) (*core.ProjectMember, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*projectpb.ProjectMember, error) {
		return c.client.GetMember(callCtx, &projectpb.GetProjectMemberRequest{Id: id})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoMember(resp)
}

func (c *Client) ListMembers(ctx context.Context, projectID string) ([]core.ProjectMember, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*projectpb.ListProjectMembersResponse, error) {
		return c.client.ListMembers(callCtx, &projectpb.ListProjectMembersRequest{ProjectId: projectID})
	})
	if err != nil {
		return nil, err
	}
	items := make([]core.ProjectMember, 0, len(resp.Members))
	for _, member := range resp.Members {
		item, convErr := fromProtoMember(member)
		if convErr != nil {
			return nil, convErr
		}
		items = append(items, *item)
	}
	return items, nil
}

func (c *Client) UpdateMember(ctx context.Context, input core.UpdateProjectMemberInput) (*core.ProjectMember, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*projectpb.ProjectMember, error) {
		return c.client.UpdateMember(callCtx, &projectpb.UpdateProjectMemberRequest{
			Id:            input.ID,
			RoleInProject: input.RoleInProject,
			IsActive:      input.IsActive,
			LeftAt:        input.LeftAt,
		})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoMember(resp)
}

func (c *Client) DeleteMember(ctx context.Context, id string) error {
	_, err := callWithTimeout(ctx, func(callCtx context.Context) (*projectpb.Empty, error) {
		return c.client.DeleteMember(callCtx, &projectpb.DeleteProjectMemberRequest{Id: id})
	})
	return err
}

func (c *Client) CreateEvent(ctx context.Context, input core.CreateProjectEventInput) (*core.ProjectEvent, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*projectpb.ProjectEvent, error) {
		return c.client.CreateEvent(callCtx, &projectpb.CreateProjectEventRequest{
			ProjectId:   input.ProjectID,
			StageId:     input.StageID,
			Name:        input.Name,
			Description: input.Description,
			PlannedDate: input.PlannedDate,
			ActualDate:  input.ActualDate,
			Status:      input.Status,
		})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoEvent(resp)
}

func (c *Client) GetEvent(ctx context.Context, id string) (*core.ProjectEvent, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*projectpb.ProjectEvent, error) {
		return c.client.GetEvent(callCtx, &projectpb.GetProjectEventRequest{Id: id})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoEvent(resp)
}

func (c *Client) ListEvents(ctx context.Context, projectID string) ([]core.ProjectEvent, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*projectpb.ListProjectEventsResponse, error) {
		return c.client.ListEvents(callCtx, &projectpb.ListProjectEventsRequest{ProjectId: projectID})
	})
	if err != nil {
		return nil, err
	}
	items := make([]core.ProjectEvent, 0, len(resp.Events))
	for _, projectEvent := range resp.Events {
		item, convErr := fromProtoEvent(projectEvent)
		if convErr != nil {
			return nil, convErr
		}
		items = append(items, *item)
	}
	return items, nil
}

func (c *Client) UpdateEvent(ctx context.Context, input core.UpdateProjectEventInput) (*core.ProjectEvent, error) {
	resp, err := callWithTimeout(ctx, func(callCtx context.Context) (*projectpb.ProjectEvent, error) {
		return c.client.UpdateEvent(callCtx, &projectpb.UpdateProjectEventRequest{
			Id:          input.ID,
			StageId:     input.StageID,
			Name:        input.Name,
			Description: input.Description,
			PlannedDate: input.PlannedDate,
			ActualDate:  input.ActualDate,
			Status:      input.Status,
		})
	})
	if err != nil {
		return nil, err
	}
	return fromProtoEvent(resp)
}

func (c *Client) DeleteEvent(ctx context.Context, id string) error {
	_, err := callWithTimeout(ctx, func(callCtx context.Context) (*projectpb.Empty, error) {
		return c.client.DeleteEvent(callCtx, &projectpb.DeleteProjectEventRequest{Id: id})
	})
	return err
}

func callWithTimeout[T any](ctx context.Context, fn func(context.Context) (T, error)) (T, error) {
	callCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return fn(callCtx)
}

func fromProtoProject(project *projectpb.Project) (*core.Project, error) {
	if project == nil {
		return nil, errors.New("project payload is empty")
	}
	return &core.Project{ID: project.ID, ProjectCode: project.ProjectCode, Name: project.Name, Description: project.Description, ContractNumber: project.ContractNumber, Status: project.Status, CustomerName: project.CustomerName, ManagerID: project.ManagerId, PlannedStartDate: project.PlannedStartDate, PlannedDeadline: project.PlannedDeadline, ActualStartDate: project.ActualStartDate, ActualDeadline: project.ActualDeadline, PlannedBudget: project.PlannedBudget, CreatedAt: project.CreatedAt, UpdatedAt: project.UpdatedAt}, nil
}

func fromProtoStage(stage *projectpb.ProjectStage) (*core.ProjectStage, error) {
	if stage == nil {
		return nil, errors.New("project stage payload is empty")
	}
	return &core.ProjectStage{ID: stage.ID, ProjectID: stage.ProjectId, Name: stage.Name, Description: stage.Description, SequenceNumber: stage.SequenceNumber, Status: stage.Status, PlannedStartDate: stage.PlannedStartDate, PlannedEndDate: stage.PlannedEndDate, ActualStartDate: stage.ActualStartDate, ActualEndDate: stage.ActualEndDate, PlannedIncome: stage.PlannedIncome, PlannedExpense: stage.PlannedExpense, CreatedAt: stage.CreatedAt, UpdatedAt: stage.UpdatedAt}, nil
}

func fromProtoMember(member *projectpb.ProjectMember) (*core.ProjectMember, error) {
	if member == nil {
		return nil, errors.New("project member payload is empty")
	}
	return &core.ProjectMember{ID: member.ID, ProjectID: member.ProjectId, UserID: member.UserId, RoleInProject: member.RoleInProject, IsActive: member.IsActive, JoinedAt: member.JoinedAt, LeftAt: member.LeftAt}, nil
}

func fromProtoEvent(projectEvent *projectpb.ProjectEvent) (*core.ProjectEvent, error) {
	if projectEvent == nil {
		return nil, errors.New("project event payload is empty")
	}
	return &core.ProjectEvent{ID: projectEvent.ID, ProjectID: projectEvent.ProjectId, StageID: projectEvent.StageId, Name: projectEvent.Name, Description: projectEvent.Description, PlannedDate: projectEvent.PlannedDate, ActualDate: projectEvent.ActualDate, Status: projectEvent.Status, CreatedAt: projectEvent.CreatedAt, UpdatedAt: projectEvent.UpdatedAt}, nil
}
