package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
)

type Service struct {
	log    *slog.Logger
	users  UserProvider
	hasher PasswordHasher
	tokens TokenManager
}

// New returns a new instance of the Auth service
func New(log *slog.Logger, users UserProvider, hasher PasswordHasher, tokens TokenManager) *Service {
	return &Service{
		log:    log,
		users:  users,
		hasher: hasher,
		tokens: tokens,
	}
}

func (s *Service) Register(ctx context.Context, name, email, password string) (User, error) {
	const op = "auth.RegisterNewUser"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email))

	log.Info("registering user")

	// TODO: Сделать валидацию в отдельную функцию
	// -----------------------------------------
	name = strings.TrimSpace(name)
	email = normalizeEmail(email)

	if name == "" || email == "" {
		return User{}, ErrInvalidArgument
	}
	if !strings.Contains(email, "@") {
		return User{}, ErrInvalidArgument
	}
	if utf8.RuneCountInString(password) < 8 {
		return User{}, ErrWeakPassword
	}
	/// ----------------------------------------

	//passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	//if err != nil {
	//	log.Error("failed to generate password hash", err)
	//	return User{}, fmt.Errorf("%s: %w", op, err)
	//}

	passHash, err := s.hasher.Hash(password)
	if err != nil {
		s.log.Error("failed to generate password hash", "error", err)
		return User{}, fmt.Errorf("%s: %w", op, err)
	}

	input := CreateUserInput{
		ID:           uuid.NewString(),
		Name:         name,
		Email:        email,
		PasswordHash: passHash,
		Role:         RoleUser,
	}

	user, err := s.users.CreateUser(ctx, input)
	if err != nil {
		s.log.Error("failed to create user", "error", err)
		return User{}, fmt.Errorf("%s: %w", op, err)
	}

	//if err != nil {
	//	//	if errors.Is(err, ErrUserExists) {
	//	//		log.Warn("user already exists", "error", err)
	//	//
	//	//		return "", fmt.Errorf("%s: %w", op, ErrUserExists)
	//	//	}
	//	//
	//	//	log.Error("failed to save user", err)
	//	//
	//	//	return "", fmt.Errorf("%s: %w", op, err)
	//	//}

	log.Info("successfully saved user")

	return user, nil

}

// Login checks if user with give credentials exists in the systems
// If user exists, but password is incorrect, returns error
// If user doesn't exist, returns error
func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	const op = "auth.Login"
	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email))

	// TODO: сделать валидатор
	email = normilizeEmail(email)
	if email == "" || password == "" {
		return "", ErrInvalidArgument
	}

	s.log.Info("attempting to login user")

	user, err := s.users.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			s.log.Warn("user not found", "error", err)
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		s.log.Error("failed to get user", "error", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err = s.hasher.Compare(user.PassHash, password); err != nil {
		s.log.Info("invalid credentials", err)
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	token, err := s.tokens.Generate(Claims{
		UserID: user.ID,
		Email:  email,
		Role:   user.Role,
	})
	if err != nil {
		s.log.Error("failed to sign token", "error", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil

	//user, err := s.db.User(ctx, email)
	//if err != nil {
	//	if errors.Is(err, ErrUserNotFound) {
	//		s.log.Warn("user not found")
	//		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	//	}
	//
	//	s.log.Error("failed to get user", "error", err)
	//
	//	return "", fmt.Errorf("%s: %w", op, err)
	//}
	//
	//if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
	//	s.log.Info("invalid credentials", err)
	//
	//	return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	//}
	//
	//app, err := s.AppProvider.App(ctx, appId)
	//if err != nil {
	//	return "", fmt.Errorf("%s: %w", op, err)
	//}
	//
	//s.log.Info("user logged in successfully")
	//
	//claims := jwt.MapClaims{
	//	"uid":    user.ID,
	//	"email":  email,
	//	"exp":    time.Now().Add(s.tokenTTL).Unix(),
	//	"app_id": appId,
	//}
	//
	//token := jwt.NewWithClaims(
	//	jwt.SigningMethodHS256,
	//	claims)
	//
	//signedToken, err := token.SignedString([]byte(app.Secret))
	//if err != nil {
	//	s.log.Error("failed to sign token", "error", err)
	//	return "", err
	//}
	//
	//return signedToken, nil
}
