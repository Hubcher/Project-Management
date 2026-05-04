package core

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IdentityService struct {
	auth  AuthClient
	users UserDirectory
}

func NewIdentityService(auth AuthClient, users UserDirectory) *IdentityService {
	return &IdentityService{auth: auth, users: users}
}

func (s *IdentityService) Register(ctx context.Context, input RegisterInput) (*RegisterResult, error) {
	userID, err := s.auth.Register(ctx, input.Email, input.Password, RoleUser)
	if err != nil {
		return nil, err
	}

	profile, err := s.users.CreateUser(ctx, CreateUserInput{
		ID:         userID,
		FirstName:  input.FirstName,
		LastName:   input.LastName,
		MiddleName: input.MiddleName,
		BirthDate:  input.BirthDate,
		Phone:      input.Phone,
		Department: input.Department,
		Position:   input.Position,
		AvatarURL:  input.AvatarURL,
		Bio:        input.Bio,
	})
	if err != nil {
		_ = s.auth.DeleteCredentials(ctx, userID)
		return nil, err
	}

	token, err := s.auth.Login(ctx, input.Email, input.Password)
	if err != nil {
		return nil, err
	}

	authUser, err := s.auth.Validate(ctx, token)
	if err != nil {
		return nil, err
	}

	return &RegisterResult{
		Token:   token,
		User:    authUser,
		Profile: *profile,
	}, nil
}

func (s *IdentityService) CreateManagedUser(ctx context.Context, input RegisterInput) (*ManagedUserResult, error) {
	role, err := validateManagedRole(input.Role)
	if err != nil {
		return nil, err
	}

	userID, err := s.auth.Register(ctx, input.Email, input.Password, role)
	if err != nil {
		return nil, err
	}

	profile, err := s.users.CreateUser(ctx, CreateUserInput{
		ID:         userID,
		FirstName:  input.FirstName,
		LastName:   input.LastName,
		MiddleName: input.MiddleName,
		BirthDate:  input.BirthDate,
		Phone:      input.Phone,
		Department: input.Department,
		Position:   input.Position,
		AvatarURL:  input.AvatarURL,
		Bio:        input.Bio,
	})
	if err != nil {
		_ = s.auth.DeleteCredentials(ctx, userID)
		return nil, err
	}

	return &ManagedUserResult{
		Email:   input.Email,
		Role:    role,
		Profile: *profile,
	}, nil
}

func (s *IdentityService) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	token, err := s.auth.Login(ctx, email, password)
	if err != nil {
		return nil, err
	}

	user, err := s.auth.Validate(ctx, token)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		Token: token,
		User:  user,
	}, nil
}

func (s *IdentityService) DeleteUser(ctx context.Context, userID string) error {
	authErr := s.auth.DeleteCredentials(ctx, userID)
	if authErr != nil && status.Code(authErr) != codes.NotFound {
		return authErr
	}

	userErr := s.users.DeleteUser(ctx, userID)
	if userErr != nil && status.Code(userErr) != codes.NotFound {
		return userErr
	}

	return nil
}

func validateManagedRole(role Role) (Role, error) {
	switch role {
	case "", RoleUser:
		return RoleUser, nil
	case RoleManager:
		return RoleManager, nil
	case RoleAdmin:
		return RoleAdmin, nil
	default:
		return "", NewStatusError(400, "invalid role")
	}
}
