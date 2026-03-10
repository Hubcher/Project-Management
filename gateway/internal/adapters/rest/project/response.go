package project

type ProjectResponse struct {
	ContractNumber int32  `json:"contract_number"`
	Name           string `json:"name"`
	StartDate      string `json:"start_date"`
	Deadline       string `json:"deadline"`
	Price          string `json:"price"`
	UserID         string `json:"user_id"`
	CreatedAt      string `json:"created_at"`
}

type ListProjectResponse struct {
	Projects []ProjectResponse `json:"projects"`
}
