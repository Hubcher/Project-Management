package project

type CreateProjectRequest struct {
	Name      string `json:"name" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	Deadline  string `json:"deadline" binding:"required"`
	Price     string `json:"price" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
}

type UpdateProjectRequest struct {
	Name      string `json:"name" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	Deadline  string `json:"deadline" binding:"required"`
	Price     string `json:"price" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
}
