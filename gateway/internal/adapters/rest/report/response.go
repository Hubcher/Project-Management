package report

type DailyReportResponse struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	ReportDate string `json:"report_date"`
	Status     string `json:"status"`
	TotalHours string `json:"total_hours"`
	Summary    string `json:"summary"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
}

type DailyReportEntryResponse struct {
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

type DailyReportCommentResponse struct {
	ID           string `json:"id"`
	ReportID     string `json:"report_id"`
	AuthorUserID string `json:"author_user_id"`
	Comment      string `json:"comment"`
	CreatedAt    string `json:"created_at,omitempty"`
}

type ListDailyReportsResponse struct{ Reports []DailyReportResponse `json:"reports"` }
type ListDailyReportEntriesResponse struct{ Entries []DailyReportEntryResponse `json:"entries"` }
type ListDailyReportCommentsResponse struct{ Comments []DailyReportCommentResponse `json:"comments"` }