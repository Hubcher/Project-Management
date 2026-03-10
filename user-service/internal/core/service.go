package core

import "log/slog"

type Service struct {
	log *slog.Logger
	db  DB
}

func NewService(log *slog.Logger, db DB) *Service {
	return &Service{
		log: log,
		db:  db,
	}
}
