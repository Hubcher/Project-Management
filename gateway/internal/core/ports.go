package core

import "context"

type Pinger interface {
	Ping(ctx context.Context) error
}

type AuthClient interface {
	Pinger
	Register(ctx context.Context, email, password string, role Role) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
	Validate(ctx context.Context, token string) (AuthUser, error)
	DeleteCredentials(ctx context.Context, userID string) error
}

type UserDirectory interface {
	Pinger
	CreateUser(ctx context.Context, input CreateUserInput) (*UserProfile, error)
	GetUser(ctx context.Context, id string) (*UserProfile, error)
	ListUsers(ctx context.Context) ([]UserProfile, error)
	UpdateUser(ctx context.Context, input UpdateUserInput) (*UserProfile, error)
	DeleteUser(ctx context.Context, id string) error
}

type ProjectDirectory interface {
	Pinger
	CreateProject(ctx context.Context, input CreateProjectInput) (*Project, error)
	GetProject(ctx context.Context, id string) (*Project, error)
	ListProjects(ctx context.Context, participantUserID string) ([]Project, error)
	UpdateProject(ctx context.Context, input UpdateProjectInput) (*Project, error)
	DeleteProject(ctx context.Context, id string) error

	CreateStage(ctx context.Context, input CreateProjectStageInput) (*ProjectStage, error)
	GetStage(ctx context.Context, id string) (*ProjectStage, error)
	ListStages(ctx context.Context, projectID string) ([]ProjectStage, error)
	UpdateStage(ctx context.Context, input UpdateProjectStageInput) (*ProjectStage, error)
	DeleteStage(ctx context.Context, id string) error

	CreateMember(ctx context.Context, input CreateProjectMemberInput) (*ProjectMember, error)
	GetMember(ctx context.Context, id string) (*ProjectMember, error)
	ListMembers(ctx context.Context, projectID string) ([]ProjectMember, error)
	UpdateMember(ctx context.Context, input UpdateProjectMemberInput) (*ProjectMember, error)
	DeleteMember(ctx context.Context, id string) error

	CreateEvent(ctx context.Context, input CreateProjectEventInput) (*ProjectEvent, error)
	GetEvent(ctx context.Context, id string) (*ProjectEvent, error)
	ListEvents(ctx context.Context, projectID string) ([]ProjectEvent, error)
	UpdateEvent(ctx context.Context, input UpdateProjectEventInput) (*ProjectEvent, error)
	DeleteEvent(ctx context.Context, id string) error
}

type ReportDirectory interface {
	Pinger
	CreateReport(ctx context.Context, input CreateDailyReportInput) (*DailyReport, error)
	GetReport(ctx context.Context, id string) (*DailyReport, error)
	ListReports(ctx context.Context, input ListDailyReportsInput) ([]DailyReport, error)
	UpdateReport(ctx context.Context, input UpdateDailyReportInput) (*DailyReport, error)
	DeleteReport(ctx context.Context, id string) error

	CreateEntry(ctx context.Context, input CreateDailyReportEntryInput) (*DailyReportEntry, error)
	GetEntry(ctx context.Context, id string) (*DailyReportEntry, error)
	ListEntries(ctx context.Context, reportID string) ([]DailyReportEntry, error)
	UpdateEntry(ctx context.Context, input UpdateDailyReportEntryInput) (*DailyReportEntry, error)
	DeleteEntry(ctx context.Context, id string) error

	CreateComment(ctx context.Context, input CreateDailyReportCommentInput) (*DailyReportComment, error)
	GetComment(ctx context.Context, id string) (*DailyReportComment, error)
	ListComments(ctx context.Context, reportID string) ([]DailyReportComment, error)
	UpdateComment(ctx context.Context, input UpdateDailyReportCommentInput) (*DailyReportComment, error)
	DeleteComment(ctx context.Context, id string) error
}