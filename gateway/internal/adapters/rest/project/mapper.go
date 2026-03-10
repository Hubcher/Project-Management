package project

import projectpb "github.com/Hubcher/project-management/contracts/gen/proto/project"

func toCreateProjectPB(req CreateProjectRequest) *projectpb.CreateProjectRequest {
	return &projectpb.CreateProjectRequest{
		Name:      req.Name,
		StartDate: req.StartDate,
		Deadline:  req.Deadline,
		Price:     req.Price,
		UserId:    req.UserID,
	}
}

func toUpdateProjectPB(contractNumber int32, req UpdateProjectRequest) *projectpb.UpdateProjectRequest {
	return &projectpb.UpdateProjectRequest{
		ContractNumber: contractNumber,
		Name:           req.Name,
		StartDate:      req.StartDate,
		Deadline:       req.Deadline,
		Price:          req.Price,
		UserId:         req.UserID,
	}
}

func toProjectResponse(p *projectpb.Project) ProjectResponse {
	return ProjectResponse{
		ContractNumber: p.GetContractNumber(),
		Name:           p.GetName(),
		StartDate:      p.GetStartDate(),
		Deadline:       p.GetDeadline(),
		Price:          p.GetPrice(),
		UserID:         p.GetUserId(),
		CreatedAt:      p.GetCreatedAt(),
	}
}

func toListProjectsResponse(items []*projectpb.Project) ListProjectResponse {
	resp := ListProjectResponse{
		Projects: make([]ProjectResponse, 0, len(items)),
	}

	for _, item := range items {
		resp.Projects = append(resp.Projects, toProjectResponse(item))
	}

	return resp
}
