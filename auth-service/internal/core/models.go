package core

type User struct {
	ID       string
	Username string
	PassHash []byte
}

type App struct {
	ID     string
	Name   string
	Secret string
}
