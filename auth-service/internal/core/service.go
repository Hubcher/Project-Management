package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	log         *slog.Logger
	db          DB
	AppProvider AppProvider
	tokenTTL    time.Duration
}

// New returns a new instance of the Auth service
func New(log *slog.Logger, db DB, provider AppProvider) *Service {
	return &Service{
		log:         log,
		db:          db,
		AppProvider: provider,
	}
}

// Login checks if user with give credentials exists in the systems
// If user exists, but password is incorrect, returns error
// If user doesn't exist, returns error
func (s *Service) Login(ctx context.Context, email string, password string, appId int) (string, error) {
	const op = "auth.Login"
	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email))

	log.Info("attempting to login user")

	user, err := s.db.User(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			s.log.Warn("user not found")
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		s.log.Error("failed to get user", "error", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		s.log.Info("invalid credentials", err)

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := s.AppProvider.App(ctx, appId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	s.log.Info("user logged in successfully")

	claims := jwt.MapClaims{
		"uid":    user.ID,
		"email":  email,
		"exp":    time.Now().Add(s.tokenTTL).Unix(),
		"app_id": appId,
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims)

	signedToken, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		s.log.Error("failed to sign token", "error", err)
		return "", err
	}

	return signedToken, nil
}

func (s *Service) RegisterNewUser(ctx context.Context, email string, password string) (string, error) {
	const op = "auth.RegisterNewUser"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email))

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	id, err := s.db.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, ErrUserExists) {
			log.Warn("user already exists", "error", err)

			return "", fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		log.Error("failed to save user", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully saved user")

	return id, nil

}

func (s *Service) IsAdmin(ctx context.Context, userID string) (bool, error) {
	const op = "auth.IsAdmin"
	log := s.log.With(
		slog.String("op", op),
		slog.String("userID", userID),
	)

	log.Info("checking if user is admin")

	isAdmin, err := s.db.isAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrAppNotFound) {
			log.Warn("user not found")
			return false, fmt.Errorf("%s: %w", op, ErrInvalidAppID)
		}
		log.Error("failed to check if user is admin", "error", err)
		return false, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("checking if user is admin", slog.Bool("isAdmin", isAdmin))

	return isAdmin, nil
}
