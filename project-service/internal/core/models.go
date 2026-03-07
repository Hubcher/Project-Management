package core

type Project struct {
	Id             int
	Name           string
	StartDate      string
	Deadline       string
	Price          string
	CreatedAt      string
	ContractNumber int
}

type NewProjectDto struct {
	Name           string
	StartDate      string
	Deadline       string
	Price          string
	ContractNumber int
}
