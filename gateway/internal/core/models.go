package core

type Role string

const (
	RoleUser    Role = "user"
	RoleManager Role = "manager"
	RoleAdmin   Role = "admin"
)

type AuthUser struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   Role   `json:"role"`
}

type UserProfile struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name,omitempty"`
	BirthDate  string `json:"birth_date,omitempty"`
	Phone      string `json:"phone,omitempty"`
	Department string `json:"department,omitempty"`
	Position   string `json:"position,omitempty"`
	AvatarURL  string `json:"avatar_url,omitempty"`
	Bio        string `json:"bio,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
}

type Project struct {
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

type ProjectStage struct {
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

type ProjectMember struct {
	ID            string `json:"id"`
	ProjectID     string `json:"project_id"`
	UserID        string `json:"user_id"`
	RoleInProject string `json:"role_in_project"`
	IsActive      bool   `json:"is_active"`
	JoinedAt      string `json:"joined_at,omitempty"`
	LeftAt        string `json:"left_at,omitempty"`
}

type ProjectEvent struct {
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

type DailyReport struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	ReportDate string `json:"report_date"`
	Status     string `json:"status"`
	TotalHours string `json:"total_hours"`
	Summary    string `json:"summary"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
}

type DailyReportEntry struct {
	ID          string `json:"id"`
	ReportID    string `json:"report_id"`
	ProjectID   string `json:"project_id"`
	StageID     string `json:"stage_id,omitempty"`
	WorkType    string `json:"work_type"`
	Description string `json:"description"`
	HoursSpent  string `json:"hours_spent"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

type DailyReportComment struct {
	ID           string `json:"id"`
	ReportID     string `json:"report_id"`
	AuthorUserID string `json:"author_user_id"`
	Comment      string `json:"comment"`
	CreatedAt    string `json:"created_at,omitempty"`
}

type Payment struct {
	ID          string `json:"id"`
	ProjectID   string `json:"project_id"`
	StageID     string `json:"stage_id,omitempty"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	PlannedDate string `json:"planned_date"`
	ActualDate  string `json:"actual_date,omitempty"`
	Description string `json:"description"`
	CreatedBy   string `json:"created_by,omitempty"`
	PaidBy      string `json:"paid_by,omitempty"`
	IsOverdue   bool   `json:"is_overdue"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

type ProjectFinancialSummary struct {
	ProjectID      string `json:"project_id"`
	PlannedIncome  string `json:"planned_income"`
	PlannedExpense string `json:"planned_expense"`
	PlannedBalance string `json:"planned_balance"`
	PaidIncome     string `json:"paid_income"`
	PaidExpense    string `json:"paid_expense"`
	PaidBalance    string `json:"paid_balance"`
	OverdueIncome  string `json:"overdue_income"`
	OverdueExpense string `json:"overdue_expense"`
	OverdueCount   int32  `json:"overdue_count"`
}

type RegisterInput struct {
	Email      string
	Password   string
	Role       Role
	FirstName  string
	LastName   string
	MiddleName string
	BirthDate  string
	Phone      string
	Department string
	Position   string
	AvatarURL  string
	Bio        string
}

type CreateUserInput struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name,omitempty"`
	BirthDate  string `json:"birth_date,omitempty"`
	Phone      string `json:"phone,omitempty"`
	Department string `json:"department,omitempty"`
	Position   string `json:"position,omitempty"`
	AvatarURL  string `json:"avatar_url,omitempty"`
	Bio        string `json:"bio,omitempty"`
}

type UpdateUserInput struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name,omitempty"`
	BirthDate  string `json:"birth_date,omitempty"`
	Phone      string `json:"phone,omitempty"`
	Department string `json:"department,omitempty"`
	Position   string `json:"position,omitempty"`
	AvatarURL  string `json:"avatar_url,omitempty"`
	Bio        string `json:"bio,omitempty"`
}

type CreateProjectInput struct {
	ProjectCode      string `json:"project_code"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	ContractNumber   string `json:"contract_number"`
	Status           string `json:"status"`
	CustomerName     string `json:"customer_name"`
	ManagerID        string `json:"manager_id,omitempty"`
	PlannedStartDate string `json:"planned_start_date,omitempty"`
	PlannedDeadline  string `json:"planned_deadline,omitempty"`
	ActualStartDate  string `json:"actual_start_date,omitempty"`
	ActualDeadline   string `json:"actual_deadline,omitempty"`
	PlannedBudget    string `json:"planned_budget"`
}

type UpdateProjectInput struct {
	ID               string `json:"id"`
	ProjectCode      string `json:"project_code"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	ContractNumber   string `json:"contract_number"`
	Status           string `json:"status"`
	CustomerName     string `json:"customer_name"`
	ManagerID        string `json:"manager_id,omitempty"`
	PlannedStartDate string `json:"planned_start_date,omitempty"`
	PlannedDeadline  string `json:"planned_deadline,omitempty"`
	ActualStartDate  string `json:"actual_start_date,omitempty"`
	ActualDeadline   string `json:"actual_deadline,omitempty"`
	PlannedBudget    string `json:"planned_budget"`
}

type CreateProjectStageInput struct {
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
}

type UpdateProjectStageInput struct {
	ID               string `json:"id"`
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
}

type CreateProjectMemberInput struct {
	ProjectID     string `json:"project_id"`
	UserID        string `json:"user_id"`
	RoleInProject string `json:"role_in_project"`
}

type UpdateProjectMemberInput struct {
	ID            string `json:"id"`
	RoleInProject string `json:"role_in_project"`
	IsActive      bool   `json:"is_active"`
	LeftAt        string `json:"left_at,omitempty"`
}

type CreateProjectEventInput struct {
	ProjectID   string `json:"project_id"`
	StageID     string `json:"stage_id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PlannedDate string `json:"planned_date,omitempty"`
	ActualDate  string `json:"actual_date,omitempty"`
	Status      string `json:"status"`
}

type UpdateProjectEventInput struct {
	ID          string `json:"id"`
	StageID     string `json:"stage_id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PlannedDate string `json:"planned_date,omitempty"`
	ActualDate  string `json:"actual_date,omitempty"`
	Status      string `json:"status"`
}

type ListDailyReportsInput struct {
	UserID   string `json:"user_id,omitempty"`
	Status   string `json:"status,omitempty"`
	DateFrom string `json:"date_from,omitempty"`
	DateTo   string `json:"date_to,omitempty"`
}

type CreateDailyReportInput struct {
	UserID     string `json:"user_id"`
	ReportDate string `json:"report_date"`
	Status     string `json:"status"`
	Summary    string `json:"summary"`
}

type UpdateDailyReportInput struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Summary string `json:"summary"`
}

type CreateDailyReportEntryInput struct {
	ReportID    string `json:"report_id"`
	ProjectID   string `json:"project_id"`
	StageID     string `json:"stage_id,omitempty"`
	WorkType    string `json:"work_type"`
	Description string `json:"description"`
	HoursSpent  string `json:"hours_spent"`
}

type UpdateDailyReportEntryInput struct {
	ID          string `json:"id"`
	ProjectID   string `json:"project_id"`
	StageID     string `json:"stage_id,omitempty"`
	WorkType    string `json:"work_type"`
	Description string `json:"description"`
	HoursSpent  string `json:"hours_spent"`
}

type CreateDailyReportCommentInput struct {
	ReportID     string `json:"report_id"`
	AuthorUserID string `json:"author_user_id"`
	Comment      string `json:"comment"`
}

type UpdateDailyReportCommentInput struct {
	ID      string `json:"id"`
	Comment string `json:"comment"`
}

type CreatePaymentInput struct {
	ProjectID   string `json:"project_id"`
	StageID     string `json:"stage_id,omitempty"`
	Type        string `json:"type"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency,omitempty"`
	PlannedDate string `json:"planned_date"`
	Description string `json:"description"`
	CreatedBy   string `json:"created_by,omitempty"`
}

type ListPaymentsInput struct {
	ProjectID   string `json:"project_id"`
	StageID     string `json:"stage_id,omitempty"`
	Type        string `json:"type,omitempty"`
	Status      string `json:"status,omitempty"`
	DateFrom    string `json:"date_from,omitempty"`
	DateTo      string `json:"date_to,omitempty"`
	OverdueOnly bool   `json:"overdue_only,omitempty"`
}

type UpdatePaymentInput struct {
	ID          string `json:"id"`
	StageID     string `json:"stage_id,omitempty"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	PlannedDate string `json:"planned_date"`
	ActualDate  string `json:"actual_date,omitempty"`
	Description string `json:"description"`
	PaidBy      string `json:"paid_by,omitempty"`
}

type MarkPaymentPaidInput struct {
	ID         string `json:"id"`
	ActualDate string `json:"actual_date,omitempty"`
	PaidBy     string `json:"paid_by,omitempty"`
}

type BuildExportInput struct {
	ReportType         string `json:"report_type"`
	Format             string `json:"format"`
	ProjectID          string `json:"project_id,omitempty"`
	DateFrom           string `json:"date_from,omitempty"`
	DateTo             string `json:"date_to,omitempty"`
	GroupBy            string `json:"group_by,omitempty"`
	PaymentType        string `json:"payment_type,omitempty"`
	PaymentStatus      string `json:"payment_status,omitempty"`
	OverdueOnly        bool   `json:"overdue_only,omitempty"`
	RequesterUserID    string `json:"requester_user_id,omitempty"`
	IncludeAllProjects bool   `json:"include_all_projects,omitempty"`
}

type ExportedFile struct {
	FileName    string
	ContentType string
	Data        []byte
}

type RegisterResult struct {
	Token   string      `json:"token"`
	User    AuthUser    `json:"user"`
	Profile UserProfile `json:"profile"`
}

type LoginResult struct {
	Token string   `json:"token"`
	User  AuthUser `json:"user"`
}

type ManagedUserResult struct {
	Email   string      `json:"email"`
	Role    Role        `json:"role"`
	Profile UserProfile `json:"profile"`
}
