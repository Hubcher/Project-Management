package core

import (
	"context"
	"strings"

	"github.com/google/uuid"
)

type service struct {
	repo ProjectRepository
}

func NewService(repo ProjectRepository) ProjectService {
	return &service{repo: repo}
}

func (s *service) CreateProject(ctx context.Context, input CreateProjectInput) (*Project, error) {
	if err := validateProjectInput(input.Name, input.StartDate, input.Deadline, input.Price, input.UserID); err != nil {
		return nil, err
	}

	return s.repo.CreateProject(ctx, input)
}

func (s *service) GetProject(ctx context.Context, contractNumber int64) (*Project, error) {
	if contractNumber <= 0 {
		return nil, ErrInvalidProject
	}

	return s.repo.GetProject(ctx, contractNumber)
}

func (s *service) ListProjects(ctx context.Context, userID string) ([]Project, error) {
	if _, err := uuid.Parse(userID); err != nil {
		return nil, ErrInvalidProject
	}

	return s.repo.ListProjects(ctx, userID)
}

func (s *service) UpdateProject(ctx context.Context, input UpdateProjectInput) (*Project, error) {
	if input.ContractNumber <= 0 {
		return nil, ErrInvalidProject
	}

	if err := validateProjectInput(input.Name, input.StartDate, input.Deadline, input.Price, input.UserID); err != nil {
		return nil, err
	}

	return s.repo.UpdateProject(ctx, input)
}

func (s *service) DeleteProject(ctx context.Context, contractNumber int64) error {
	if contractNumber <= 0 {
		return ErrInvalidProject
	}

	return s.repo.DeleteProject(ctx, contractNumber)
}

func validateProjectInput(name string, startDate, deadline interface{ Before(interface{}) bool }, price, userID string) error {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(price) == "" {
		return ErrInvalidProject
	}

	if _, err := uuid.Parse(userID); err != nil {
		return ErrInvalidProject
	}

	return nil
}
