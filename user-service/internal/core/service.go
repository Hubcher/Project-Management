package core

import (
    "context"
    "net/url"
    "strings"
    "time"

    "github.com/google/uuid"
)

type Service struct {
    db DB
}

func NewService(db DB) *Service {
    return &Service{db: db}
}

func (s *Service) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
    input = normalizeCreateInput(input)
    if err := validateProfile(input.ID, input.FirstName, input.LastName, input.BirthDate, input.Phone, input.AvatarURL); err != nil {
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

func (s *Service) ListUsers(ctx context.Context) ([]User, error) {
    return s.db.ListUsers(ctx)
}

func (s *Service) UpdateUser(ctx context.Context, input UpdateUserInput) (*User, error) {
    input = normalizeUpdateInput(input)
    if err := validateProfile(input.ID, input.FirstName, input.LastName, input.BirthDate, input.Phone, input.AvatarURL); err != nil {
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

func normalizeCreateInput(input CreateUserInput) CreateUserInput {
    input.FirstName = strings.TrimSpace(input.FirstName)
    input.LastName = strings.TrimSpace(input.LastName)
    input.MiddleName = strings.TrimSpace(input.MiddleName)
    input.Phone = strings.TrimSpace(input.Phone)
    input.Department = strings.TrimSpace(input.Department)
    input.Position = strings.TrimSpace(input.Position)
    input.AvatarURL = strings.TrimSpace(input.AvatarURL)
    input.Bio = strings.TrimSpace(input.Bio)
    if input.BirthDate != nil {
        day := input.BirthDate.UTC()
        normalized := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
        input.BirthDate = &normalized
    }
    return input
}

func normalizeUpdateInput(input UpdateUserInput) UpdateUserInput {
    input.FirstName = strings.TrimSpace(input.FirstName)
    input.LastName = strings.TrimSpace(input.LastName)
    input.MiddleName = strings.TrimSpace(input.MiddleName)
    input.Phone = strings.TrimSpace(input.Phone)
    input.Department = strings.TrimSpace(input.Department)
    input.Position = strings.TrimSpace(input.Position)
    input.AvatarURL = strings.TrimSpace(input.AvatarURL)
    input.Bio = strings.TrimSpace(input.Bio)
    if input.BirthDate != nil {
        day := input.BirthDate.UTC()
        normalized := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
        input.BirthDate = &normalized
    }
    return input
}

func validateProfile(id, firstName, lastName string, birthDate *time.Time, phone, avatarURL string) error {
    if !isValidUUID(id) {
        return ErrInvalidUser
    }
    if firstName == "" || lastName == "" {
        return ErrInvalidUser
    }
    if birthDate != nil && birthDate.After(time.Now()) {
        return ErrInvalidUser
    }
    if len(phone) > 32 {
        return ErrInvalidUser
    }
    if avatarURL != "" {
        if _, err := url.ParseRequestURI(avatarURL); err != nil {
            return ErrInvalidUser
        }
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
        strings.Contains(msg, "users_pkey")
}
