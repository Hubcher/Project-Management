package core

import (
    "context"
    "errors"
    "log/slog"
    "strings"

    "github.com/google/uuid"
)

type Service struct {
    log  *slog.Logger
    repo AuthRepository
    tm   TokenManager
    pm   PasswordManager
}

func NewService(log *slog.Logger, repo AuthRepository, tm TokenManager, pm PasswordManager) *Service {
    return &Service{
        log:  log,
        repo: repo,
        tm:   tm,
        pm:   pm,
    }
}

func (s *Service) Register(ctx context.Context, email, password string) (string, error) {
    return s.createAccount(ctx, email, password, RoleUser)
}

func (s *Service) RegisterWithRole(ctx context.Context, email, password string, role Role) (string, error) {
    return s.createAccount(ctx, email, password, role)
}

func (s *Service) EnsureAdmin(ctx context.Context, email, password string) error {
    const op = "auth.Service.EnsureAdmin"
    log := s.log.With(slog.String("op", op), slog.String("email", email))

    email = strings.TrimSpace(strings.ToLower(email))
    password = strings.TrimSpace(password)
    if email == "" || password == "" {
        return ErrInvalidArgument
    }

    acc, err := s.repo.GetByEmail(ctx, email)
    if err == nil {
        if acc.Role != RoleAdmin {
            log.Error("bootstrap admin email already exists with non-admin role")
            return ErrForbidden
        }
        log.Info("bootstrap admin already exists, skipping")
        return nil
    }
    if !errors.Is(err, ErrAccountNotFound) {
        return err
    }

    _, err = s.createAccount(ctx, email, password, RoleAdmin)
    if err != nil {
        return err
    }

    log.Info("bootstrap admin created")
    return nil
}

func (s *Service) createAccount(ctx context.Context, email, password string, role Role) (string, error) {
    const op = "auth.Service.createAccount"
    log := s.log.With(
        slog.String("op", op),
        slog.String("email", email),
        slog.String("role", string(role)),
    )

    email = strings.TrimSpace(strings.ToLower(email))
    password = strings.TrimSpace(password)
    role = normalizeRole(role)

    if email == "" || password == "" {
        return "", ErrInvalidArgument
    }
    if role != RoleUser && role != RoleAdmin {
        return "", ErrInvalidArgument
    }

    if _, err := s.repo.GetByEmail(ctx, email); err == nil {
        log.Info("user with this email already exists")
        return "", ErrAccountExists
    } else if !errors.Is(err, ErrAccountNotFound) {
        log.Error("failed to check existing account", "error", err)
        return "", err
    }

    passHash, err := s.pm.Hash(password)
    if err != nil {
        log.Error("failed to generate password hash", "error", err)
        return "", err
    }

    id := uuid.NewString()
    acc := Account{
        UserID:       id,
        Email:        email,
        PasswordHash: passHash,
        Role:         role,
        IsActive:     true,
    }

    if err = s.repo.CreateAccount(ctx, acc); err != nil {
        return "", err
    }

    return id, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
    const op = "auth.Service.Login"
    log := s.log.With(
        slog.String("op", op),
        slog.String("email", email),
    )

    log.Info("checking user")

    email = strings.TrimSpace(strings.ToLower(email))
    password = strings.TrimSpace(password)

    if email == "" || password == "" {
        return "", ErrInvalidArgument
    }

    acc, err := s.repo.GetByEmail(ctx, email)
    if err != nil {
        log.Error("failed to get user by email", "error", err)
        return "", ErrInvalidCredentials
    }

    if !acc.IsActive {
        return "", ErrInactiveAccount
    }

    if err = s.pm.Compare(acc.PasswordHash, password); err != nil {
        log.Error("failed to compare password", "error", err)
        return "", ErrInvalidCredentials
    }

    return s.tm.Generate(Claims{
        UserID: acc.UserID,
        Role:   acc.Role,
        Email:  acc.Email,
    })
}

func (s *Service) Validate(ctx context.Context, token string) (Claims, error) {
    token = strings.TrimSpace(token)
    if token == "" {
        return Claims{}, ErrInvalidToken
    }

    claims, err := s.tm.Parse(token)
    if err != nil {
        return Claims{}, err
    }

    acc, err := s.repo.GetByUserID(ctx, claims.UserID)
    if err != nil {
        return Claims{}, ErrInvalidToken
    }

    if !acc.IsActive {
        return Claims{}, ErrInactiveAccount
    }

    return Claims{
        UserID: acc.UserID,
        Role:   acc.Role,
        Email:  acc.Email,
    }, nil
}

func (s *Service) DeleteCredentials(ctx context.Context, userID string) error {
    userID = strings.TrimSpace(userID)
    if userID == "" {
        return ErrInvalidArgument
    }
    return s.repo.DeleteByUserID(ctx, userID)
}

func normalizeRole(role Role) Role {
    switch strings.ToLower(strings.TrimSpace(string(role))) {
    case string(RoleAdmin):
        return RoleAdmin
    default:
        return RoleUser
    }
}
