package projectpb

import (
	"encoding/json"

	"google.golang.org/grpc/encoding"
)

const JSONCodecName = "json"

type jsonCodec struct{}

func (jsonCodec) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (jsonCodec) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (jsonCodec) Name() string {
	return JSONCodecName
}

func init() {
	encoding.RegisterCodec(jsonCodec{})
}

type Empty struct{}

type Project struct {
	ID               string `json:"id,omitempty"`
	ProjectCode      string `json:"project_code,omitempty"`
	Name             string `json:"name,omitempty"`
	Description      string `json:"description,omitempty"`
	ContractNumber   string `json:"contract_number,omitempty"`
	Status           string `json:"status,omitempty"`
	CustomerName     string `json:"customer_name,omitempty"`
	ManagerId        string `json:"manager_id,omitempty"`
	PlannedStartDate string `json:"planned_start_date,omitempty"`
	PlannedDeadline  string `json:"planned_deadline,omitempty"`
	ActualStartDate  string `json:"actual_start_date,omitempty"`
	ActualDeadline   string `json:"actual_deadline,omitempty"`
	PlannedBudget    string `json:"planned_budget,omitempty"`
	CreatedAt        string `json:"created_at,omitempty"`
	UpdatedAt        string `json:"updated_at,omitempty"`
}

type CreateProjectRequest struct {
	ProjectCode      string `json:"project_code,omitempty"`
	Name             string `json:"name,omitempty"`
	Description      string `json:"description,omitempty"`
	ContractNumber   string `json:"contract_number,omitempty"`
	Status           string `json:"status,omitempty"`
	CustomerName     string `json:"customer_name,omitempty"`
	ManagerId        string `json:"manager_id,omitempty"`
	PlannedStartDate string `json:"planned_start_date,omitempty"`
	PlannedDeadline  string `json:"planned_deadline,omitempty"`
	ActualStartDate  string `json:"actual_start_date,omitempty"`
	ActualDeadline   string `json:"actual_deadline,omitempty"`
	PlannedBudget    string `json:"planned_budget,omitempty"`
}

type GetProjectRequest struct {
	Id string `json:"id,omitempty"`
}

type ListProjectsRequest struct {
	ParticipantUserId string `json:"participant_user_id,omitempty"`
}

type ListProjectsResponse struct {
	Projects []*Project `json:"projects,omitempty"`
}

type UpdateProjectRequest struct {
	Id               string `json:"id,omitempty"`
	ProjectCode      string `json:"project_code,omitempty"`
	Name             string `json:"name,omitempty"`
	Description      string `json:"description,omitempty"`
	ContractNumber   string `json:"contract_number,omitempty"`
	Status           string `json:"status,omitempty"`
	CustomerName     string `json:"customer_name,omitempty"`
	ManagerId        string `json:"manager_id,omitempty"`
	PlannedStartDate string `json:"planned_start_date,omitempty"`
	PlannedDeadline  string `json:"planned_deadline,omitempty"`
	ActualStartDate  string `json:"actual_start_date,omitempty"`
	ActualDeadline   string `json:"actual_deadline,omitempty"`
	PlannedBudget    string `json:"planned_budget,omitempty"`
}

type DeleteProjectRequest struct {
	Id string `json:"id,omitempty"`
}

type ProjectStage struct {
	ID               string `json:"id,omitempty"`
	ProjectId        string `json:"project_id,omitempty"`
	Name             string `json:"name,omitempty"`
	Description      string `json:"description,omitempty"`
	SequenceNumber   int32  `json:"sequence_number,omitempty"`
	Status           string `json:"status,omitempty"`
	PlannedStartDate string `json:"planned_start_date,omitempty"`
	PlannedEndDate   string `json:"planned_end_date,omitempty"`
	ActualStartDate  string `json:"actual_start_date,omitempty"`
	ActualEndDate    string `json:"actual_end_date,omitempty"`
	PlannedIncome    string `json:"planned_income,omitempty"`
	PlannedExpense   string `json:"planned_expense,omitempty"`
	CreatedAt        string `json:"created_at,omitempty"`
	UpdatedAt        string `json:"updated_at,omitempty"`
}

type CreateProjectStageRequest struct {
	ProjectId        string `json:"project_id,omitempty"`
	Name             string `json:"name,omitempty"`
	Description      string `json:"description,omitempty"`
	SequenceNumber   int32  `json:"sequence_number,omitempty"`
	Status           string `json:"status,omitempty"`
	PlannedStartDate string `json:"planned_start_date,omitempty"`
	PlannedEndDate   string `json:"planned_end_date,omitempty"`
	ActualStartDate  string `json:"actual_start_date,omitempty"`
	ActualEndDate    string `json:"actual_end_date,omitempty"`
	PlannedIncome    string `json:"planned_income,omitempty"`
	PlannedExpense   string `json:"planned_expense,omitempty"`
}

type GetProjectStageRequest struct {
	Id string `json:"id,omitempty"`
}

type ListProjectStagesRequest struct {
	ProjectId string `json:"project_id,omitempty"`
}

type ListProjectStagesResponse struct {
	Stages []*ProjectStage `json:"stages,omitempty"`
}

type UpdateProjectStageRequest struct {
	Id               string `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	Description      string `json:"description,omitempty"`
	SequenceNumber   int32  `json:"sequence_number,omitempty"`
	Status           string `json:"status,omitempty"`
	PlannedStartDate string `json:"planned_start_date,omitempty"`
	PlannedEndDate   string `json:"planned_end_date,omitempty"`
	ActualStartDate  string `json:"actual_start_date,omitempty"`
	ActualEndDate    string `json:"actual_end_date,omitempty"`
	PlannedIncome    string `json:"planned_income,omitempty"`
	PlannedExpense   string `json:"planned_expense,omitempty"`
}

type DeleteProjectStageRequest struct {
	Id string `json:"id,omitempty"`
}

type ProjectMember struct {
	ID            string `json:"id,omitempty"`
	ProjectId     string `json:"project_id,omitempty"`
	UserId        string `json:"user_id,omitempty"`
	RoleInProject string `json:"role_in_project,omitempty"`
	IsActive      bool   `json:"is_active"`
	JoinedAt      string `json:"joined_at,omitempty"`
	LeftAt        string `json:"left_at,omitempty"`
}

type CreateProjectMemberRequest struct {
	ProjectId     string `json:"project_id,omitempty"`
	UserId        string `json:"user_id,omitempty"`
	RoleInProject string `json:"role_in_project,omitempty"`
}

type GetProjectMemberRequest struct {
	Id string `json:"id,omitempty"`
}

type ListProjectMembersRequest struct {
	ProjectId string `json:"project_id,omitempty"`
}

type ListProjectMembersResponse struct {
	Members []*ProjectMember `json:"members,omitempty"`
}

type UpdateProjectMemberRequest struct {
	Id            string `json:"id,omitempty"`
	RoleInProject string `json:"role_in_project,omitempty"`
	IsActive      bool   `json:"is_active"`
	LeftAt        string `json:"left_at,omitempty"`
}

type DeleteProjectMemberRequest struct {
	Id string `json:"id,omitempty"`
}

type ProjectEvent struct {
	ID          string `json:"id,omitempty"`
	ProjectId   string `json:"project_id,omitempty"`
	StageId     string `json:"stage_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	PlannedDate string `json:"planned_date,omitempty"`
	ActualDate  string `json:"actual_date,omitempty"`
	Status      string `json:"status,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

type CreateProjectEventRequest struct {
	ProjectId   string `json:"project_id,omitempty"`
	StageId     string `json:"stage_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	PlannedDate string `json:"planned_date,omitempty"`
	ActualDate  string `json:"actual_date,omitempty"`
	Status      string `json:"status,omitempty"`
}

type GetProjectEventRequest struct {
	Id string `json:"id,omitempty"`
}

type ListProjectEventsRequest struct {
	ProjectId string `json:"project_id,omitempty"`
}

type ListProjectEventsResponse struct {
	Events []*ProjectEvent `json:"events,omitempty"`
}

type UpdateProjectEventRequest struct {
	Id          string `json:"id,omitempty"`
	StageId     string `json:"stage_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	PlannedDate string `json:"planned_date,omitempty"`
	ActualDate  string `json:"actual_date,omitempty"`
	Status      string `json:"status,omitempty"`
}

type DeleteProjectEventRequest struct {
	Id string `json:"id,omitempty"`
}
