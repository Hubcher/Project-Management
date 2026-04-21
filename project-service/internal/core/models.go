package core

import "time"

type ProjectStatus string

type StageStatus string

type EventStatus string

const (
	ProjectStatusDraft     ProjectStatus = "draft"
	ProjectStatusPlanned   ProjectStatus = "planned"
	ProjectStatusActive    ProjectStatus = "active"
	ProjectStatusPaused    ProjectStatus = "paused"
	ProjectStatusCompleted ProjectStatus = "completed"
	ProjectStatusCancelled ProjectStatus = "cancelled"

	StageStatusDraft      StageStatus = "draft"
	StageStatusPlanned    StageStatus = "planned"
	StageStatusInProgress StageStatus = "in_progress"
	StageStatusCompleted  StageStatus = "completed"
	StageStatusCancelled  StageStatus = "cancelled"

	EventStatusPlanned   EventStatus = "planned"
	EventStatusReached   EventStatus = "reached"
	EventStatusCancelled EventStatus = "cancelled"
)

type Project struct {
	ID               string
	ProjectCode      string
	Name             string
	Description      string
	ContractNumber   string
	Status           ProjectStatus
	CustomerName     string
	ManagerID        string
	PlannedStartDate *time.Time
	PlannedDeadline  *time.Time
	ActualStartDate  *time.Time
	ActualDeadline   *time.Time
	PlannedBudget    string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type CreateProjectInput struct {
	ProjectCode      string
	Name             string
	Description      string
	ContractNumber   string
	Status           ProjectStatus
	CustomerName     string
	ManagerID        string
	PlannedStartDate *time.Time
	PlannedDeadline  *time.Time
	ActualStartDate  *time.Time
	ActualDeadline   *time.Time
	PlannedBudget    string
}

type UpdateProjectInput struct {
	ID               string
	ProjectCode      string
	Name             string
	Description      string
	ContractNumber   string
	Status           ProjectStatus
	CustomerName     string
	ManagerID        string
	PlannedStartDate *time.Time
	PlannedDeadline  *time.Time
	ActualStartDate  *time.Time
	ActualDeadline   *time.Time
	PlannedBudget    string
}

type ProjectStage struct {
	ID               string
	ProjectID        string
	Name             string
	Description      string
	SequenceNumber   int32
	Status           StageStatus
	PlannedStartDate *time.Time
	PlannedEndDate   *time.Time
	ActualStartDate  *time.Time
	ActualEndDate    *time.Time
	PlannedIncome    string
	PlannedExpense   string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type CreateProjectStageInput struct {
	ProjectID        string
	Name             string
	Description      string
	SequenceNumber   int32
	Status           StageStatus
	PlannedStartDate *time.Time
	PlannedEndDate   *time.Time
	ActualStartDate  *time.Time
	ActualEndDate    *time.Time
	PlannedIncome    string
	PlannedExpense   string
}

type UpdateProjectStageInput struct {
	ID               string
	Name             string
	Description      string
	SequenceNumber   int32
	Status           StageStatus
	PlannedStartDate *time.Time
	PlannedEndDate   *time.Time
	ActualStartDate  *time.Time
	ActualEndDate    *time.Time
	PlannedIncome    string
	PlannedExpense   string
}

type ProjectMember struct {
	ID            string
	ProjectID     string
	UserID        string
	RoleInProject string
	IsActive      bool
	JoinedAt      time.Time
	LeftAt        *time.Time
}

type CreateProjectMemberInput struct {
	ProjectID     string
	UserID        string
	RoleInProject string
}

type UpdateProjectMemberInput struct {
	ID            string
	RoleInProject string
	IsActive      bool
	LeftAt        *time.Time
}

type ProjectEvent struct {
	ID          string
	ProjectID   string
	StageID     string
	Name        string
	Description string
	PlannedDate *time.Time
	ActualDate  *time.Time
	Status      EventStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateProjectEventInput struct {
	ProjectID   string
	StageID     string
	Name        string
	Description string
	PlannedDate *time.Time
	ActualDate  *time.Time
	Status      EventStatus
}

type UpdateProjectEventInput struct {
	ID          string
	StageID     string
	Name        string
	Description string
	PlannedDate *time.Time
	ActualDate  *time.Time
	Status      EventStatus
}
