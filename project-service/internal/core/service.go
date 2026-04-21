package core

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"
)

type service struct {
	repo ProjectRepository
}

func NewService(repo ProjectRepository) ProjectService {
	return &service{repo: repo}
}

func (s *service) CreateProject(ctx context.Context, input CreateProjectInput) (*Project, error) {
	input = normalizeCreateProjectInput(input)
	if err := validateProjectInput(input.ProjectCode, input.Name, input.ContractNumber, input.ManagerID, input.Status, input.PlannedStartDate, input.PlannedDeadline, input.ActualStartDate, input.ActualDeadline, input.PlannedBudget); err != nil {
		return nil, err
	}
	return s.repo.CreateProject(ctx, input)
}

func (s *service) GetProject(ctx context.Context, id string) (*Project, error) {
	id = strings.TrimSpace(id)
	if !isValidUUID(id) {
		return nil, ErrInvalidProject
	}
	return s.repo.GetProject(ctx, id)
}

func (s *service) ListProjects(ctx context.Context, participantUserID string) ([]Project, error) {
	participantUserID = strings.TrimSpace(participantUserID)
	if participantUserID != "" && !isValidUUID(participantUserID) {
		return nil, ErrInvalidProject
	}
	return s.repo.ListProjects(ctx, participantUserID)
}

func (s *service) UpdateProject(ctx context.Context, input UpdateProjectInput) (*Project, error) {
	input = normalizeUpdateProjectInput(input)
	if !isValidUUID(input.ID) {
		return nil, ErrInvalidProject
	}
	if err := validateProjectInput(input.ProjectCode, input.Name, input.ContractNumber, input.ManagerID, input.Status, input.PlannedStartDate, input.PlannedDeadline, input.ActualStartDate, input.ActualDeadline, input.PlannedBudget); err != nil {
		return nil, err
	}
	return s.repo.UpdateProject(ctx, input)
}

func (s *service) DeleteProject(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if !isValidUUID(id) {
		return ErrInvalidProject
	}
	return s.repo.DeleteProject(ctx, id)
}

func (s *service) CreateStage(ctx context.Context, input CreateProjectStageInput) (*ProjectStage, error) {
	input = normalizeCreateStageInput(input)
	if !isValidUUID(input.ProjectID) {
		return nil, ErrInvalidStage
	}
	if _, err := s.repo.GetProject(ctx, input.ProjectID); err != nil {
		return nil, err
	}
	if err := validateStageInput(input.Name, input.SequenceNumber, input.Status, input.PlannedStartDate, input.PlannedEndDate, input.ActualStartDate, input.ActualEndDate, input.PlannedIncome, input.PlannedExpense); err != nil {
		return nil, err
	}
	return s.repo.CreateStage(ctx, input)
}

func (s *service) GetStage(ctx context.Context, id string) (*ProjectStage, error) {
	id = strings.TrimSpace(id)
	if !isValidUUID(id) {
		return nil, ErrInvalidStage
	}
	return s.repo.GetStage(ctx, id)
}

func (s *service) ListStages(ctx context.Context, projectID string) ([]ProjectStage, error) {
	projectID = strings.TrimSpace(projectID)
	if !isValidUUID(projectID) {
		return nil, ErrInvalidStage
	}
	if _, err := s.repo.GetProject(ctx, projectID); err != nil {
		return nil, err
	}
	return s.repo.ListStages(ctx, projectID)
}

func (s *service) UpdateStage(ctx context.Context, input UpdateProjectStageInput) (*ProjectStage, error) {
	input = normalizeUpdateStageInput(input)
	if !isValidUUID(input.ID) {
		return nil, ErrInvalidStage
	}
	if err := validateStageInput(input.Name, input.SequenceNumber, input.Status, input.PlannedStartDate, input.PlannedEndDate, input.ActualStartDate, input.ActualEndDate, input.PlannedIncome, input.PlannedExpense); err != nil {
		return nil, err
	}
	return s.repo.UpdateStage(ctx, input)
}

func (s *service) DeleteStage(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if !isValidUUID(id) {
		return ErrInvalidStage
	}
	return s.repo.DeleteStage(ctx, id)
}

func (s *service) CreateMember(ctx context.Context, input CreateProjectMemberInput) (*ProjectMember, error) {
	input = normalizeCreateMemberInput(input)
	if !isValidUUID(input.ProjectID) || !isValidUUID(input.UserID) {
		return nil, ErrInvalidMember
	}
	if _, err := s.repo.GetProject(ctx, input.ProjectID); err != nil {
		return nil, err
	}
	if strings.TrimSpace(input.RoleInProject) == "" {
		return nil, ErrInvalidMember
	}
	return s.repo.CreateMember(ctx, input)
}

func (s *service) GetMember(ctx context.Context, id string) (*ProjectMember, error) {
	id = strings.TrimSpace(id)
	if !isValidUUID(id) {
		return nil, ErrInvalidMember
	}
	return s.repo.GetMember(ctx, id)
}

func (s *service) ListMembers(ctx context.Context, projectID string) ([]ProjectMember, error) {
	projectID = strings.TrimSpace(projectID)
	if !isValidUUID(projectID) {
		return nil, ErrInvalidMember
	}
	if _, err := s.repo.GetProject(ctx, projectID); err != nil {
		return nil, err
	}
	return s.repo.ListMembers(ctx, projectID)
}

func (s *service) UpdateMember(ctx context.Context, input UpdateProjectMemberInput) (*ProjectMember, error) {
	input = normalizeUpdateMemberInput(input)
	if !isValidUUID(input.ID) || strings.TrimSpace(input.RoleInProject) == "" {
		return nil, ErrInvalidMember
	}
	if input.IsActive {
		input.LeftAt = nil
	} else if input.LeftAt == nil {
		now := time.Now().UTC()
		input.LeftAt = &now
	}
	return s.repo.UpdateMember(ctx, input)
}

func (s *service) DeleteMember(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if !isValidUUID(id) {
		return ErrInvalidMember
	}
	return s.repo.DeleteMember(ctx, id)
}

func (s *service) CreateEvent(ctx context.Context, input CreateProjectEventInput) (*ProjectEvent, error) {
	input = normalizeCreateEventInput(input)
	if !isValidUUID(input.ProjectID) {
		return nil, ErrInvalidEvent
	}
	if _, err := s.repo.GetProject(ctx, input.ProjectID); err != nil {
		return nil, err
	}
	if err := validateEventInput(input.ProjectID, input.StageID, input.Name, input.Status); err != nil {
		return nil, err
	}
	if err := s.ensureStageBelongsToProject(ctx, input.StageID, input.ProjectID); err != nil {
		return nil, err
	}
	return s.repo.CreateEvent(ctx, input)
}

func (s *service) GetEvent(ctx context.Context, id string) (*ProjectEvent, error) {
	id = strings.TrimSpace(id)
	if !isValidUUID(id) {
		return nil, ErrInvalidEvent
	}
	return s.repo.GetEvent(ctx, id)
}

func (s *service) ListEvents(ctx context.Context, projectID string) ([]ProjectEvent, error) {
	projectID = strings.TrimSpace(projectID)
	if !isValidUUID(projectID) {
		return nil, ErrInvalidEvent
	}
	if _, err := s.repo.GetProject(ctx, projectID); err != nil {
		return nil, err
	}
	return s.repo.ListEvents(ctx, projectID)
}

func (s *service) UpdateEvent(ctx context.Context, input UpdateProjectEventInput) (*ProjectEvent, error) {
	input = normalizeUpdateEventInput(input)
	if !isValidUUID(input.ID) {
		return nil, ErrInvalidEvent
	}
	current, err := s.repo.GetEvent(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if err = validateEventInput(current.ProjectID, input.StageID, input.Name, input.Status); err != nil {
		return nil, err
	}
	if err = s.ensureStageBelongsToProject(ctx, input.StageID, current.ProjectID); err != nil {
		return nil, err
	}
	return s.repo.UpdateEvent(ctx, input)
}

func (s *service) DeleteEvent(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if !isValidUUID(id) {
		return ErrInvalidEvent
	}
	return s.repo.DeleteEvent(ctx, id)
}

func (s *service) ensureStageBelongsToProject(ctx context.Context, stageID, projectID string) error {
	if strings.TrimSpace(stageID) == "" {
		return nil
	}
	stage, err := s.repo.GetStage(ctx, stageID)
	if err != nil {
		return err
	}
	if stage.ProjectID != projectID {
		return ErrInvalidEvent
	}
	return nil
}

func validateProjectInput(projectCode, name, contractNumber, managerID string, status ProjectStatus, plannedStart, plannedDeadline, actualStart, actualDeadline *time.Time, plannedBudget string) error {
	if strings.TrimSpace(projectCode) == "" || strings.TrimSpace(name) == "" || !isValidUUID(contractNumber) || !isValidUUID(managerID) {
		return ErrInvalidProject
	}
	if !isValidProjectStatus(status) || !isValidMoney(plannedBudget) {
		return ErrInvalidProject
	}
	if !isValidDateRange(plannedStart, plannedDeadline) || !isValidDateRange(actualStart, actualDeadline) {
		return ErrInvalidProject
	}
	return nil
}

func validateStageInput(name string, sequenceNumber int32, status StageStatus, plannedStart, plannedEnd, actualStart, actualEnd *time.Time, plannedIncome, plannedExpense string) error {
	if strings.TrimSpace(name) == "" || sequenceNumber <= 0 {
		return ErrInvalidStage
	}
	if !isValidStageStatus(status) || !isValidMoney(plannedIncome) || !isValidMoney(plannedExpense) {
		return ErrInvalidStage
	}
	if !isValidDateRange(plannedStart, plannedEnd) || !isValidDateRange(actualStart, actualEnd) {
		return ErrInvalidStage
	}
	return nil
}

func validateEventInput(projectID, stageID, name string, status EventStatus) error {
	if !isValidUUID(projectID) || strings.TrimSpace(name) == "" {
		return ErrInvalidEvent
	}
	if strings.TrimSpace(stageID) != "" && !isValidUUID(stageID) {
		return ErrInvalidEvent
	}
	if !isValidEventStatus(status) {
		return ErrInvalidEvent
	}
	return nil
}

func normalizeCreateProjectInput(input CreateProjectInput) CreateProjectInput {
	input.ProjectCode = strings.TrimSpace(input.ProjectCode)
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)
	input.ContractNumber = normalizeUUIDOrGenerate(input.ContractNumber)
	input.CustomerName = strings.TrimSpace(input.CustomerName)
	input.ManagerID = strings.TrimSpace(input.ManagerID)
	input.Status = normalizeProjectStatus(input.Status, true)
	input.PlannedBudget = normalizeMoney(input.PlannedBudget)
	input.PlannedStartDate = normalizeDate(input.PlannedStartDate)
	input.PlannedDeadline = normalizeDate(input.PlannedDeadline)
	input.ActualStartDate = normalizeDate(input.ActualStartDate)
	input.ActualDeadline = normalizeDate(input.ActualDeadline)
	return input
}

func normalizeUpdateProjectInput(input UpdateProjectInput) UpdateProjectInput {
	input.ID = strings.TrimSpace(input.ID)
	input.ProjectCode = strings.TrimSpace(input.ProjectCode)
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)
	input.ContractNumber = strings.TrimSpace(input.ContractNumber)
	input.CustomerName = strings.TrimSpace(input.CustomerName)
	input.ManagerID = strings.TrimSpace(input.ManagerID)
	input.Status = normalizeProjectStatus(input.Status, false)
	input.PlannedBudget = normalizeMoney(input.PlannedBudget)
	input.PlannedStartDate = normalizeDate(input.PlannedStartDate)
	input.PlannedDeadline = normalizeDate(input.PlannedDeadline)
	input.ActualStartDate = normalizeDate(input.ActualStartDate)
	input.ActualDeadline = normalizeDate(input.ActualDeadline)
	return input
}

func normalizeCreateStageInput(input CreateProjectStageInput) CreateProjectStageInput {
	input.ProjectID = strings.TrimSpace(input.ProjectID)
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)
	input.Status = normalizeStageStatus(input.Status, true)
	input.PlannedIncome = normalizeMoney(input.PlannedIncome)
	input.PlannedExpense = normalizeMoney(input.PlannedExpense)
	input.PlannedStartDate = normalizeDate(input.PlannedStartDate)
	input.PlannedEndDate = normalizeDate(input.PlannedEndDate)
	input.ActualStartDate = normalizeDate(input.ActualStartDate)
	input.ActualEndDate = normalizeDate(input.ActualEndDate)
	return input
}

func normalizeUpdateStageInput(input UpdateProjectStageInput) UpdateProjectStageInput {
	input.ID = strings.TrimSpace(input.ID)
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)
	input.Status = normalizeStageStatus(input.Status, false)
	input.PlannedIncome = normalizeMoney(input.PlannedIncome)
	input.PlannedExpense = normalizeMoney(input.PlannedExpense)
	input.PlannedStartDate = normalizeDate(input.PlannedStartDate)
	input.PlannedEndDate = normalizeDate(input.PlannedEndDate)
	input.ActualStartDate = normalizeDate(input.ActualStartDate)
	input.ActualEndDate = normalizeDate(input.ActualEndDate)
	return input
}

func normalizeCreateMemberInput(input CreateProjectMemberInput) CreateProjectMemberInput {
	input.ProjectID = strings.TrimSpace(input.ProjectID)
	input.UserID = strings.TrimSpace(input.UserID)
	role := strings.TrimSpace(input.RoleInProject)
	if role == "" {
		role = "member"
	}
	input.RoleInProject = role
	return input
}

func normalizeUpdateMemberInput(input UpdateProjectMemberInput) UpdateProjectMemberInput {
	input.ID = strings.TrimSpace(input.ID)
	input.RoleInProject = strings.TrimSpace(input.RoleInProject)
	input.LeftAt = normalizeTimestamp(input.LeftAt)
	return input
}

func normalizeCreateEventInput(input CreateProjectEventInput) CreateProjectEventInput {
	input.ProjectID = strings.TrimSpace(input.ProjectID)
	input.StageID = strings.TrimSpace(input.StageID)
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)
	input.Status = normalizeEventStatus(input.Status, true)
	input.PlannedDate = normalizeDate(input.PlannedDate)
	input.ActualDate = normalizeDate(input.ActualDate)
	return input
}

func normalizeUpdateEventInput(input UpdateProjectEventInput) UpdateProjectEventInput {
	input.ID = strings.TrimSpace(input.ID)
	input.StageID = strings.TrimSpace(input.StageID)
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)
	input.Status = normalizeEventStatus(input.Status, false)
	input.PlannedDate = normalizeDate(input.PlannedDate)
	input.ActualDate = normalizeDate(input.ActualDate)
	return input
}

func normalizeUUIDOrGenerate(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return uuid.NewString()
	}
	return value
}

func normalizeMoney(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "0"
	}
	return value
}

func normalizeProjectStatus(status ProjectStatus, allowDefault bool) ProjectStatus {
	status = ProjectStatus(strings.ToLower(strings.TrimSpace(string(status))))
	if status == "" && allowDefault {
		return ProjectStatusDraft
	}
	return status
}

func normalizeStageStatus(status StageStatus, allowDefault bool) StageStatus {
	status = StageStatus(strings.ToLower(strings.TrimSpace(string(status))))
	if status == "" && allowDefault {
		return StageStatusDraft
	}
	return status
}

func normalizeEventStatus(status EventStatus, allowDefault bool) EventStatus {
	status = EventStatus(strings.ToLower(strings.TrimSpace(string(status))))
	if status == "" && allowDefault {
		return EventStatusPlanned
	}
	return status
}

func normalizeDate(value *time.Time) *time.Time {
	if value == nil || value.IsZero() {
		return nil
	}
	day := value.UTC()
	normalized := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	return &normalized
}

func normalizeTimestamp(value *time.Time) *time.Time {
	if value == nil || value.IsZero() {
		return nil
	}
	ts := value.UTC()
	return &ts
}

func isValidUUID(value string) bool {
	_, err := uuid.Parse(strings.TrimSpace(value))
	return err == nil
}

func isValidMoney(value string) bool {
	var amount big.Rat
	if _, ok := amount.SetString(strings.TrimSpace(value)); !ok {
		return false
	}
	return amount.Sign() >= 0
}

func isValidDateRange(start, end *time.Time) bool {
	if start == nil || end == nil {
		return true
	}
	return !end.Before(*start)
}

func isValidProjectStatus(status ProjectStatus) bool {
	switch status {
	case ProjectStatusDraft, ProjectStatusPlanned, ProjectStatusActive, ProjectStatusPaused, ProjectStatusCompleted, ProjectStatusCancelled:
		return true
	default:
		return false
	}
}

func isValidStageStatus(status StageStatus) bool {
	switch status {
	case StageStatusDraft, StageStatusPlanned, StageStatusInProgress, StageStatusCompleted, StageStatusCancelled:
		return true
	default:
		return false
	}
}

func isValidEventStatus(status EventStatus) bool {
	switch status {
	case EventStatusPlanned, EventStatusReached, EventStatusCancelled:
		return true
	default:
		return false
	}
}

func IsNotFound(err error) bool {
	return errors.Is(err, ErrProjectNotFound) || errors.Is(err, ErrStageNotFound) || errors.Is(err, ErrMemberNotFound) || errors.Is(err, ErrEventNotFound)
}
