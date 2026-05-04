package exportpb

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

type BuildExportRequest struct {
	ReportType         string `json:"report_type,omitempty"`
	Format             string `json:"format,omitempty"`
	ProjectId          string `json:"project_id,omitempty"`
	DateFrom           string `json:"date_from,omitempty"`
	DateTo             string `json:"date_to,omitempty"`
	GroupBy            string `json:"group_by,omitempty"`
	PaymentType        string `json:"payment_type,omitempty"`
	PaymentStatus      string `json:"payment_status,omitempty"`
	OverdueOnly        bool   `json:"overdue_only,omitempty"`
	RequesterUserId    string `json:"requester_user_id,omitempty"`
	IncludeAllProjects bool   `json:"include_all_projects,omitempty"`
}

type BuildExportResponse struct {
	FileName    string `json:"file_name,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Data        []byte `json:"data,omitempty"`
}
