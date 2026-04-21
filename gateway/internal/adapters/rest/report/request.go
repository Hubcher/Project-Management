package report

type CreateDailyReportRequest struct {
	UserID     string `json:"user_id"`
	ReportDate string `json:"report_date"`
	Status     string `json:"status"`
	Summary    string `json:"summary"`
}

type UpdateDailyReportRequest struct {
	Status  *string `json:"status"`
	Summary *string `json:"summary"`
}

type CreateDailyReportEntryRequest struct {
	ProjectID   string `json:"project_id"`
	StageID     string `json:"stage_id"`
	WorkType    string `json:"work_type"`
	Description string `json:"description"`
	HoursSpent  string `json:"hours_spent"`
}

type UpdateDailyReportEntryRequest struct {
	ProjectID   *string `json:"project_id"`
	StageID     *string `json:"stage_id"`
	WorkType    *string `json:"work_type"`
	Description *string `json:"description"`
	HoursSpent  *string `json:"hours_spent"`
}

type CreateDailyReportCommentRequest struct {
	Comment string `json:"comment"`
}

type UpdateDailyReportCommentRequest struct {
	Comment string `json:"comment"`
}