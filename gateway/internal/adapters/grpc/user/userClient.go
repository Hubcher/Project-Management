package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	userpb "github.com/Hubcher/project-management/contracts/gen/go/user"
	"github.com/Hubcher/project-management/gateway/internal/core"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	log    *slog.Logger
	client userpb.UserServiceClient
	conn   *grpc.ClientConn
}

type profileEnvelope struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name,omitempty"`
	BirthDate  string `json:"birth_date,omitempty"`
	Phone      string `json:"phone,omitempty"`
	Department string `json:"department,omitempty"`
	Position   string `json:"position,omitempty"`
	AvatarURL  string `json:"avatar_url,omitempty"`
	Bio        string `json:"bio,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  time.Second,
				Multiplier: 1.6,
				MaxDelay:   5 * time.Second,
			},
			MinConnectTimeout: 10 * time.Second,
		}),
	)
	if err != nil {
		return nil, err
	}
	conn.Connect()

	return &Client{
		client: userpb.NewUserServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.client.Ping(ctx, &emptypb.Empty{})
	return err
}

func (c *Client) CreateUser(ctx context.Context, input core.CreateUserInput) (*core.UserProfile, error) {
	payload, err := encodeProfile(input.FirstName, input.LastName, input.MiddleName, input.BirthDate, input.Phone, input.Department, input.Position, input.AvatarURL, input.Bio, "")
	if err != nil {
		return nil, err
	}

	resp, err := c.client.CreateUser(ctx, &userpb.CreateUserRequest{
		Id:   input.ID,
		Name: payload,
	})
	if err != nil {
		return nil, err
	}
	return fromProtoUser(resp)
}

func (c *Client) GetUser(ctx context.Context, id string) (*core.UserProfile, error) {
	resp, err := c.client.GetUserById(ctx, &userpb.GetUserByIdRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return fromProtoUser(resp)
}

func (c *Client) ListUsers(ctx context.Context) ([]core.UserProfile, error) {
	resp, err := c.client.ListUsers(ctx, &userpb.ListUsersRequest{})
	if err != nil {
		return nil, err
	}

	users := make([]core.UserProfile, 0, len(resp.GetUsers()))
	for _, user := range resp.GetUsers() {
		profile, convErr := fromProtoUser(user)
		if convErr != nil {
			return nil, convErr
		}
		users = append(users, *profile)
	}
	return users, nil
}

func (c *Client) UpdateUser(ctx context.Context, input core.UpdateUserInput) (*core.UserProfile, error) {
	payload, err := encodeProfile(input.FirstName, input.LastName, input.MiddleName, input.BirthDate, input.Phone, input.Department, input.Position, input.AvatarURL, input.Bio, "")
	if err != nil {
		return nil, err
	}

	resp, err := c.client.UpdateUser(ctx, &userpb.UpdateUserRequest{
		Id:   input.ID,
		Name: &wrappers.StringValue{Value: payload},
	})
	if err != nil {
		return nil, err
	}
	return fromProtoUser(resp)
}

func (c *Client) DeleteUser(ctx context.Context, id string) error {
	_, err := c.client.DeleteUser(ctx, &userpb.DeleteUserRequest{Id: id})
	return err
}

func encodeProfile(firstName, lastName, middleName, birthDate, phone, department, position, avatarURL, bio, updatedAt string) (string, error) {
	payload, err := json.Marshal(profileEnvelope{
		FirstName:  firstName,
		LastName:   lastName,
		MiddleName: middleName,
		BirthDate:  birthDate,
		Phone:      phone,
		Department: department,
		Position:   position,
		AvatarURL:  avatarURL,
		Bio:        bio,
		UpdatedAt:  updatedAt,
	})
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

func fromProtoUser(user *userpb.User) (*core.UserProfile, error) {
	if user == nil {
		return nil, errors.New("user payload is empty")
	}

	var envelope profileEnvelope
	if err := json.Unmarshal([]byte(user.GetName()), &envelope); err != nil {
		return nil, fmt.Errorf("decode user profile: %w", err)
	}

	profile := &core.UserProfile{
		ID:         user.GetId(),
		FirstName:  envelope.FirstName,
		LastName:   envelope.LastName,
		MiddleName: envelope.MiddleName,
		BirthDate:  envelope.BirthDate,
		Phone:      envelope.Phone,
		Department: envelope.Department,
		Position:   envelope.Position,
		AvatarURL:  envelope.AvatarURL,
		Bio:        envelope.Bio,
		UpdatedAt:  envelope.UpdatedAt,
	}

	if createdAt := timestampToRFC3339(user.GetCreatedAt()); createdAt != "" {
		profile.CreatedAt = createdAt
	}
	return profile, nil
}

func timestampToRFC3339(ts *timestamp.Timestamp) string {
	if ts == nil {
		return ""
	}
	return time.Unix(ts.Seconds, int64(ts.Nanos)).UTC().Format(time.RFC3339)
}
