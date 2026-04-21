package project

import "github.com/Hubcher/project-management/gateway/internal/core"

func toCreateProjectInput(req CreateProjectRequest, managerID string) core.CreateProjectInput {
	return core.CreateProjectInput{ProjectCode: req.ProjectCode, Name: req.Name, Description: req.Description, ContractNumber: req.ContractNumber, Status: req.Status, CustomerName: req.CustomerName, ManagerID: managerID, PlannedStartDate: req.PlannedStartDate, PlannedDeadline: req.PlannedDeadline, ActualStartDate: req.ActualStartDate, ActualDeadline: req.ActualDeadline, PlannedBudget: req.PlannedBudget}
}

func toUpdateProjectInput(id string, req UpdateProjectRequest, managerID string) core.UpdateProjectInput {
	return core.UpdateProjectInput{ID: id, ProjectCode: req.ProjectCode, Name: req.Name, Description: req.Description, ContractNumber: req.ContractNumber, Status: req.Status, CustomerName: req.CustomerName, ManagerID: managerID, PlannedStartDate: req.PlannedStartDate, PlannedDeadline: req.PlannedDeadline, ActualStartDate: req.ActualStartDate, ActualDeadline: req.ActualDeadline, PlannedBudget: req.PlannedBudget}
}

func toCreateStageInput(projectID string, req CreateProjectStageRequest) core.CreateProjectStageInput {
	return core.CreateProjectStageInput{ProjectID: projectID, Name: req.Name, Description: req.Description, SequenceNumber: req.SequenceNumber, Status: req.Status, PlannedStartDate: req.PlannedStartDate, PlannedEndDate: req.PlannedEndDate, ActualStartDate: req.ActualStartDate, ActualEndDate: req.ActualEndDate, PlannedIncome: req.PlannedIncome, PlannedExpense: req.PlannedExpense}
}

func toUpdateStageInput(id string, req UpdateProjectStageRequest) core.UpdateProjectStageInput {
	return core.UpdateProjectStageInput{ID: id, Name: req.Name, Description: req.Description, SequenceNumber: req.SequenceNumber, Status: req.Status, PlannedStartDate: req.PlannedStartDate, PlannedEndDate: req.PlannedEndDate, ActualStartDate: req.ActualStartDate, ActualEndDate: req.ActualEndDate, PlannedIncome: req.PlannedIncome, PlannedExpense: req.PlannedExpense}
}

func toCreateMemberInput(projectID string, req CreateProjectMemberRequest) core.CreateProjectMemberInput {
	return core.CreateProjectMemberInput{ProjectID: projectID, UserID: req.UserID, RoleInProject: req.RoleInProject}
}

func toUpdateMemberInput(id string, req UpdateProjectMemberRequest) core.UpdateProjectMemberInput {
	return core.UpdateProjectMemberInput{ID: id, RoleInProject: req.RoleInProject, IsActive: req.IsActive, LeftAt: req.LeftAt}
}

func toCreateEventInput(projectID string, req CreateProjectEventRequest) core.CreateProjectEventInput {
	return core.CreateProjectEventInput{ProjectID: projectID, StageID: req.StageID, Name: req.Name, Description: req.Description, PlannedDate: req.PlannedDate, ActualDate: req.ActualDate, Status: req.Status}
}

func toUpdateEventInput(id string, req UpdateProjectEventRequest) core.UpdateProjectEventInput {
	return core.UpdateProjectEventInput{ID: id, StageID: req.StageID, Name: req.Name, Description: req.Description, PlannedDate: req.PlannedDate, ActualDate: req.ActualDate, Status: req.Status}
}

func toProjectResponse(project *core.Project) ProjectResponse {
	return ProjectResponse{ID: project.ID, ProjectCode: project.ProjectCode, Name: project.Name, Description: project.Description, ContractNumber: project.ContractNumber, Status: project.Status, CustomerName: project.CustomerName, ManagerID: project.ManagerID, PlannedStartDate: project.PlannedStartDate, PlannedDeadline: project.PlannedDeadline, ActualStartDate: project.ActualStartDate, ActualDeadline: project.ActualDeadline, PlannedBudget: project.PlannedBudget, CreatedAt: project.CreatedAt, UpdatedAt: project.UpdatedAt}
}

func toStageResponse(stage *core.ProjectStage) ProjectStageResponse {
	return ProjectStageResponse{ID: stage.ID, ProjectID: stage.ProjectID, Name: stage.Name, Description: stage.Description, SequenceNumber: stage.SequenceNumber, Status: stage.Status, PlannedStartDate: stage.PlannedStartDate, PlannedEndDate: stage.PlannedEndDate, ActualStartDate: stage.ActualStartDate, ActualEndDate: stage.ActualEndDate, PlannedIncome: stage.PlannedIncome, PlannedExpense: stage.PlannedExpense, CreatedAt: stage.CreatedAt, UpdatedAt: stage.UpdatedAt}
}

func toMemberResponse(member *core.ProjectMember) ProjectMemberResponse {
	return ProjectMemberResponse{ID: member.ID, ProjectID: member.ProjectID, UserID: member.UserID, RoleInProject: member.RoleInProject, IsActive: member.IsActive, JoinedAt: member.JoinedAt, LeftAt: member.LeftAt}
}

func toEventResponse(projectEvent *core.ProjectEvent) ProjectEventResponse {
	return ProjectEventResponse{ID: projectEvent.ID, ProjectID: projectEvent.ProjectID, StageID: projectEvent.StageID, Name: projectEvent.Name, Description: projectEvent.Description, PlannedDate: projectEvent.PlannedDate, ActualDate: projectEvent.ActualDate, Status: projectEvent.Status, CreatedAt: projectEvent.CreatedAt, UpdatedAt: projectEvent.UpdatedAt}
}

func toListProjectsResponse(projects []core.Project) ListProjectsResponse {
	resp := ListProjectsResponse{Projects: make([]ProjectResponse, 0, len(projects))}
	for i := range projects { project := projects[i]; resp.Projects = append(resp.Projects, toProjectResponse(&project)) }
	return resp
}

func toListStagesResponse(stages []core.ProjectStage) ListProjectStagesResponse {
	resp := ListProjectStagesResponse{Stages: make([]ProjectStageResponse, 0, len(stages))}
	for i := range stages { stage := stages[i]; resp.Stages = append(resp.Stages, toStageResponse(&stage)) }
	return resp
}

func toListMembersResponse(members []core.ProjectMember) ListProjectMembersResponse {
	resp := ListProjectMembersResponse{Members: make([]ProjectMemberResponse, 0, len(members))}
	for i := range members { member := members[i]; resp.Members = append(resp.Members, toMemberResponse(&member)) }
	return resp
}

func toListEventsResponse(events []core.ProjectEvent) ListProjectEventsResponse {
	resp := ListProjectEventsResponse{Events: make([]ProjectEventResponse, 0, len(events))}
	for i := range events { projectEvent := events[i]; resp.Events = append(resp.Events, toEventResponse(&projectEvent)) }
	return resp
}
