package core

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type Service struct {
	log *slog.Logger
	db  DB
}

// добавить DB в аргуементы
func NewService(log *slog.Logger, storage sqlx.DB) *Service {
	return &Service{
		log: log,
		db:  db,
	}
}

//func (s *Service) CreateProject(ctx context.Context, dto NewProjectDto) (*Project, error) {
//	// Должны будем обработать ошибки
//	// 1) проект уже существует
//	// 2) Проверить на обязательные поля dto (напирмер, название: "Project name cannot be empty")
//	// 3) Сделать dto, которое будем возвращать
//
//}
//
//func (s *Service) GetAllProjects(ctx context.Context) ([]Project, error) {}
//
//func (s *Service) GetProjectById(ctx context.Context, id int) (*Project, error) {}
//
//func (s *Service) GetProjectByName(ctx context.Context, name string) (*Project, error) {}
//
//func (s *Service) UpdateProject(ctx context.Context, project *Project) (*Project, error) {}
//
//func (s *Service) DeleteProject(ctx context.Context, id int) (int, error) {}
