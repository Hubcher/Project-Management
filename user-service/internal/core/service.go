package core

import (
	"context"
	"log/slog"
	"net/mail"
	"strings"

	"github.com/google/uuid"
)

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

func (s *Service) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
	if err := validateUserInput(input.ID, input.Name, input.Password, input.Email, input.Role); err != nil {
		return nil, err
	}

	user, err := s.db.CreateUser(ctx, input)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}
	return user, nil
}

func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
	if !isValidUUID(id) {
		return nil, ErrInvalidUser
	}

	return s.db.GetUser(ctx, id)
}

func (s *Service) ListUsers(ctx context.Context, role string) ([]User, error) {
	if role != "" && strings.TrimSpace(role) == "" {
		return nil, ErrInvalidUser
	}

	return s.db.ListUsers(ctx, role)
}

func (s *Service) UpdateUser(ctx context.Context, input UpdateUserInput) (*User, error) {
	if err := validateUserInput(input.ID, input.Name, input.Password, input.Email, input.Role); err != nil {
		return nil, err
	}

	user, err := s.db.UpdateUser(ctx, input)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrUserAlreadyExists
		}
		return nil, err
	}

	return user, nil
}

func (s *Service) DeleteUser(ctx context.Context, id string) error {
	if !isValidUUID(id) {
		return ErrInvalidUser
	}

	return s.db.DeleteUser(ctx, id)
}

func validateUserInput(id, name, email, password, role string) error {

	if !isValidUUID(id) {
		return ErrInvalidUser
	}

	if strings.TrimSpace(name) == "" {
		return ErrInvalidUser
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return ErrInvalidUser
	}

	if strings.TrimSpace(password) == "" {
		return ErrInvalidUser
	}

	if strings.TrimSpace(role) == "" {
		return ErrInvalidUser
	}

	return nil
}

func isValidUUID(value string) bool {
	_, err := uuid.Parse(value)
	return err == nil
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "unique constraint") ||
		strings.Contains(msg, "users_email_key")
}
