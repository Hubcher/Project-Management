package grpc

import (
    "context"
    "encoding/json"
    "errors"
    "time"

    userpb "github.com/Hubcher/project-management/contracts/gen/go/user"
    "github.com/Hubcher/project-management/user-service/internal/core"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    wrappers "github.com/golang/protobuf/ptypes/wrappers"
    timestamp "github.com/golang/protobuf/ptypes/timestamp"
    "google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
    userpb.UnimplementedUserServiceServer
    service core.UserService
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

func NewServer(service core.UserService) *Server {
    return &Server{service: service}
}

func (s *Server) Ping(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
    return &emptypb.Empty{}, nil
}

func (s *Server) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.User, error) {
    profile, err := decodeProfile(req.GetName())
    if err != nil {
        return nil, status.Error(codes.InvalidArgument, "invalid profile payload")
    }

    birthDate, err := parseBirthDate(profile.BirthDate)
    if err != nil {
        return nil, status.Error(codes.InvalidArgument, "invalid birth_date")
    }

    user, err := s.service.CreateUser(ctx, core.CreateUserInput{
        ID:         req.GetId(),
        FirstName:  profile.FirstName,
        LastName:   profile.LastName,
        MiddleName: profile.MiddleName,
        BirthDate:  birthDate,
        Phone:      profile.Phone,
        Department: profile.Department,
        Position:   profile.Position,
        AvatarURL:  profile.AvatarURL,
        Bio:        profile.Bio,
    })
    if err != nil {
        return nil, mapCoreError(err)
    }

    return toProtoUser(user)
}

func (s *Server) GetUserById(ctx context.Context, req *userpb.GetUserByIdRequest) (*userpb.User, error) {
    user, err := s.service.GetUser(ctx, req.GetId())
    if err != nil {
        return nil, mapCoreError(err)
    }
    return toProtoUser(user)
}

func (s *Server) ListUsers(ctx context.Context, _ *userpb.ListUsersRequest) (*userpb.ListUsersResponse, error) {
    users, err := s.service.ListUsers(ctx)
    if err != nil {
        return nil, mapCoreError(err)
    }

    resp := &userpb.ListUsersResponse{Users: make([]*userpb.User, 0, len(users))}
    for i := range users {
        protoUser, convErr := toProtoUser(&users[i])
        if convErr != nil {
            return nil, status.Error(codes.Internal, "failed to encode user profile")
        }
        resp.Users = append(resp.Users, protoUser)
    }
    return resp, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.User, error) {
    nameValue := req.GetName()
    if nameValue == nil {
        return nil, status.Error(codes.InvalidArgument, "profile payload is required")
    }

    profile, err := decodeProfile(nameValue.GetValue())
    if err != nil {
        return nil, status.Error(codes.InvalidArgument, "invalid profile payload")
    }

    birthDate, err := parseBirthDate(profile.BirthDate)
    if err != nil {
        return nil, status.Error(codes.InvalidArgument, "invalid birth_date")
    }

    user, err := s.service.UpdateUser(ctx, core.UpdateUserInput{
        ID:         req.GetId(),
        FirstName:  profile.FirstName,
        LastName:   profile.LastName,
        MiddleName: profile.MiddleName,
        BirthDate:  birthDate,
        Phone:      profile.Phone,
        Department: profile.Department,
        Position:   profile.Position,
        AvatarURL:  profile.AvatarURL,
        Bio:        profile.Bio,
    })
    if err != nil {
        return nil, mapCoreError(err)
    }

    return toProtoUser(user)
}

func (s *Server) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*emptypb.Empty, error) {
    if err := s.service.DeleteUser(ctx, req.GetId()); err != nil {
        return nil, mapCoreError(err)
    }
    return &emptypb.Empty{}, nil
}

func toProtoUser(user *core.User) (*userpb.User, error) {
    payload, err := encodeProfile(user)
    if err != nil {
        return nil, err
    }

    return &userpb.User{
        Id:   user.ID,
        Name: payload,
        CreatedAt: &timestamp.Timestamp{
            Seconds: user.CreatedAt.Unix(),
            Nanos:   int32(user.CreatedAt.Nanosecond()),
        },
    }, nil
}

func encodeProfile(user *core.User) (string, error) {
    profile := profileEnvelope{
        FirstName:  user.FirstName,
        LastName:   user.LastName,
        MiddleName: user.MiddleName,
        Phone:      user.Phone,
        Department: user.Department,
        Position:   user.Position,
        AvatarURL:  user.AvatarURL,
        Bio:        user.Bio,
        UpdatedAt:  user.UpdatedAt.Format(time.RFC3339),
    }
    if user.BirthDate != nil {
        profile.BirthDate = user.BirthDate.Format("2006-01-02")
    }

    data, err := json.Marshal(profile)
    if err != nil {
        return "", err
    }
    return string(data), nil
}

func decodeProfile(value string) (profileEnvelope, error) {
    if value == "" {
        return profileEnvelope{}, errors.New("empty profile payload")
    }

    var profile profileEnvelope
    if err := json.Unmarshal([]byte(value), &profile); err != nil {
        return profileEnvelope{}, err
    }
    return profile, nil
}

func parseBirthDate(value string) (*time.Time, error) {
    if value == "" {
        return nil, nil
    }

    parsed, err := time.Parse("2006-01-02", value)
    if err != nil {
        return nil, err
    }

    day := parsed.UTC()
    normalized := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
    return &normalized, nil
}

func stringValue(value string) *wrappers.StringValue {
    return &wrappers.StringValue{Value: value}
}

func mapCoreError(err error) error {
    switch {
    case errors.Is(err, core.ErrInvalidUser):
        return status.Error(codes.InvalidArgument, err.Error())
    case errors.Is(err, core.ErrUserNotFound):
        return status.Error(codes.NotFound, err.Error())
    case errors.Is(err, core.ErrUserAlreadyExists):
        return status.Error(codes.AlreadyExists, err.Error())
    default:
        return status.Error(codes.Internal, "internal error")
    }
}
