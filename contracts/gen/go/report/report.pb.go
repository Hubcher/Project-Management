package reportpb

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

type DailyReport struct {
	ID         string `json:"id,omitempty"`
	UserId     string `json:"user_id,omitempty"`
	ReportDate string `json:"report_date,omitempty"`
	Status     string `json:"status,omitempty"`
	TotalHours string `json:"total_hours,omitempty"`
	Summary    string `json:"summary,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
}

type CreateDailyReportRequest struct {
	UserId     string `json:"user_id,omitempty"`
	ReportDate string `json:"report_date,omitempty"`
	Status     string `json:"status,omitempty"`
	Summary    string `json:"summary,omitempty"`
}

type GetDailyReportRequest struct {
	Id string `json:"id,omitempty"`
}

type ListDailyReportsRequest struct {
	UserId   string `json:"user_id,omitempty"`
	Status   string `json:"status,omitempty"`
	DateFrom string `json:"date_from,omitempty"`
	DateTo   string `json:"date_to,omitempty"`
}

type ListDailyReportsResponse struct {
	Reports []*DailyReport `json:"reports,omitempty"`
}

type UpdateDailyReportRequest struct {
	Id      string `json:"id,omitempty"`
	Status  string `json:"status,omitempty"`
	Summary string `json:"summary,omitempty"`
}

type DeleteDailyReportRequest struct {
	Id string `json:"id,omitempty"`
}

type DailyReportEntry struct {
	ID          string `json:"id,omitempty"`
	ReportId    string `json:"report_id,omitempty"`
	ProjectId   string `json:"project_id,omitempty"`
	StageId     string `json:"stage_id,omitempty"`
	WorkType    string `json:"work_type,omitempty"`
	Description string `json:"description,omitempty"`
	HoursSpent  string `json:"hours_spent,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

type CreateDailyReportEntryRequest struct {
	ReportId    string `json:"report_id,omitempty"`
	ProjectId   string `json:"project_id,omitempty"`
	StageId     string `json:"stage_id,omitempty"`
	WorkType    string `json:"work_type,omitempty"`
	Description string `json:"description,omitempty"`
	HoursSpent  string `json:"hours_spent,omitempty"`
}

type GetDailyReportEntryRequest struct {
	Id string `json:"id,omitempty"`
}

type ListDailyReportEntriesRequest struct {
	ReportId string `json:"report_id,omitempty"`
}

type ListDailyReportEntriesResponse struct {
	Entries []*DailyReportEntry `json:"entries,omitempty"`
}

type UpdateDailyReportEntryRequest struct {
	Id          string `json:"id,omitempty"`
	ProjectId   string `json:"project_id,omitempty"`
	StageId     string `json:"stage_id,omitempty"`
	WorkType    string `json:"work_type,omitempty"`
	Description string `json:"description,omitempty"`
	HoursSpent  string `json:"hours_spent,omitempty"`
}

type DeleteDailyReportEntryRequest struct {
	Id string `json:"id,omitempty"`
}

type DailyReportComment struct {
	ID           string `json:"id,omitempty"`
	ReportId     string `json:"report_id,omitempty"`
	AuthorUserId string `json:"author_user_id,omitempty"`
	Comment      string `json:"comment,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
}

type CreateDailyReportCommentRequest struct {
	ReportId     string `json:"report_id,omitempty"`
	AuthorUserId string `json:"author_user_id,omitempty"`
	Comment      string `json:"comment,omitempty"`
}

type GetDailyReportCommentRequest struct {
	Id string `json:"id,omitempty"`
}

type ListDailyReportCommentsRequest struct {
	ReportId string `json:"report_id,omitempty"`
}

type ListDailyReportCommentsResponse struct {
	Comments []*DailyReportComment `json:"comments,omitempty"`
}

type UpdateDailyReportCommentRequest struct {
	Id      string `json:"id,omitempty"`
	Comment string `json:"comment,omitempty"`
}

type DeleteDailyReportCommentRequest struct {
	Id string `json:"id,omitempty"`
}
