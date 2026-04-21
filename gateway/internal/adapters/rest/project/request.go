package project

type CreateProjectRequest struct {
	ProjectCode      string `json:"project_code"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	ContractNumber   string `json:"contract_number"`
	Status           string `json:"status"`
	CustomerName     string `json:"customer_name"`
	ManagerID        string `json:"manager_id"`
	PlannedStartDate string `json:"planned_start_date"`
	PlannedDeadline  string `json:"planned_deadline"`
	ActualStartDate  string `json:"actual_start_date"`
	ActualDeadline   string `json:"actual_deadline"`
	PlannedBudget    string `json:"planned_budget"`
}

type UpdateProjectRequest struct {
	ProjectCode      string `json:"project_code"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	ContractNumber   string `json:"contract_number"`
	Status           string `json:"status"`
	CustomerName     string `json:"customer_name"`
	ManagerID        string `json:"manager_id"`
	PlannedStartDate string `json:"planned_start_date"`
	PlannedDeadline  string `json:"planned_deadline"`
	ActualStartDate  string `json:"actual_start_date"`
	ActualDeadline   string `json:"actual_deadline"`
	PlannedBudget    string `json:"planned_budget"`
}

type CreateProjectStageRequest struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	SequenceNumber   int32  `json:"sequence_number"`
	Status           string `json:"status"`
	PlannedStartDate string `json:"planned_start_date"`
	PlannedEndDate   string `json:"planned_end_date"`
	ActualStartDate  string `json:"actual_start_date"`
	ActualEndDate    string `json:"actual_end_date"`
	PlannedIncome    string `json:"planned_income"`
	PlannedExpense   string `json:"planned_expense"`
}

type UpdateProjectStageRequest struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	SequenceNumber   int32  `json:"sequence_number"`
	Status           string `json:"status"`
	PlannedStartDate string `json:"planned_start_date"`
	PlannedEndDate   string `json:"planned_end_date"`
	ActualStartDate  string `json:"actual_start_date"`
	ActualEndDate    string `json:"actual_end_date"`
	PlannedIncome    string `json:"planned_income"`
	PlannedExpense   string `json:"planned_expense"`
}

type CreateProjectMemberRequest struct {
	UserID        string `json:"user_id"`
	RoleInProject string `json:"role_in_project"`
}

type UpdateProjectMemberRequest struct {
	RoleInProject string `json:"role_in_project"`
	IsActive      bool   `json:"is_active"`
	LeftAt        string `json:"left_at"`
}

type CreateProjectEventRequest struct {
	StageID     string `json:"stage_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PlannedDate string `json:"planned_date"`
	ActualDate  string `json:"actual_date"`
	Status      string `json:"status"`
}

type UpdateProjectEventRequest struct {
	StageID     string `json:"stage_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PlannedDate string `json:"planned_date"`
	ActualDate  string `json:"actual_date"`
	Status      string `json:"status"`
}
