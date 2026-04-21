package project

type ProjectResponse struct {
	ID               string `json:"id"`
	ProjectCode      string `json:"project_code"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	ContractNumber   string `json:"contract_number"`
	Status           string `json:"status"`
	CustomerName     string `json:"customer_name"`
	ManagerID        string `json:"manager_id"`
	PlannedStartDate string `json:"planned_start_date,omitempty"`
	PlannedDeadline  string `json:"planned_deadline,omitempty"`
	ActualStartDate  string `json:"actual_start_date,omitempty"`
	ActualDeadline   string `json:"actual_deadline,omitempty"`
	PlannedBudget    string `json:"planned_budget"`
	CreatedAt        string `json:"created_at,omitempty"`
	UpdatedAt        string `json:"updated_at,omitempty"`
}

type ProjectStageResponse struct {
	ID               string `json:"id"`
	ProjectID        string `json:"project_id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	SequenceNumber   int32  `json:"sequence_number"`
	Status           string `json:"status"`
	PlannedStartDate string `json:"planned_start_date,omitempty"`
	PlannedEndDate   string `json:"planned_end_date,omitempty"`
	ActualStartDate  string `json:"actual_start_date,omitempty"`
	ActualEndDate    string `json:"actual_end_date,omitempty"`
	PlannedIncome    string `json:"planned_income"`
	PlannedExpense   string `json:"planned_expense"`
	CreatedAt        string `json:"created_at,omitempty"`
	UpdatedAt        string `json:"updated_at,omitempty"`
}

type ProjectMemberResponse struct {
	ID            string `json:"id"`
	ProjectID     string `json:"project_id"`
	UserID        string `json:"user_id"`
	RoleInProject string `json:"role_in_project"`
	IsActive      bool   `json:"is_active"`
	JoinedAt      string `json:"joined_at,omitempty"`
	LeftAt        string `json:"left_at,omitempty"`
}

type ProjectEventResponse struct {
	ID          string `json:"id"`
	ProjectID   string `json:"project_id"`
	StageID     string `json:"stage_id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PlannedDate string `json:"planned_date,omitempty"`
	ActualDate  string `json:"actual_date,omitempty"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

type ListProjectsResponse struct{ Projects []ProjectResponse `json:"projects"` }
type ListProjectStagesResponse struct{ Stages []ProjectStageResponse `json:"stages"` }
type ListProjectMembersResponse struct{ Members []ProjectMemberResponse `json:"members"` }
type ListProjectEventsResponse struct{ Events []ProjectEventResponse `json:"events"` }
