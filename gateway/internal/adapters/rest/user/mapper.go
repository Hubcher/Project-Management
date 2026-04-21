package user

import (
    "strings"

    "github.com/Hubcher/project-management/gateway/internal/core"
)

func toRegisterInput(req CreateUserRequest) core.RegisterInput {
    return core.RegisterInput{
        Email:      req.Email,
        Password:   req.Password,
        Role:       core.Role(strings.ToLower(strings.TrimSpace(req.Role))),
        FirstName:  req.FirstName,
        LastName:   req.LastName,
        MiddleName: req.MiddleName,
        BirthDate:  req.BirthDate,
        Phone:      req.Phone,
        Department: req.Department,
        Position:   req.Position,
        AvatarURL:  req.AvatarURL,
        Bio:        req.Bio,
    }
}

func toUpdateInput(id string, req UpdateUserRequest) core.UpdateUserInput {
    return core.UpdateUserInput{
        ID:         id,
        FirstName:  req.FirstName,
        LastName:   req.LastName,
        MiddleName: req.MiddleName,
        BirthDate:  req.BirthDate,
        Phone:      req.Phone,
        Department: req.Department,
        Position:   req.Position,
        AvatarURL:  req.AvatarURL,
        Bio:        req.Bio,
    }
}

func toUserResponse(profile *core.UserProfile) UserResponse {
    return UserResponse{
        ID:         profile.ID,
        FirstName:  profile.FirstName,
        LastName:   profile.LastName,
        MiddleName: profile.MiddleName,
        BirthDate:  profile.BirthDate,
        Phone:      profile.Phone,
        Department: profile.Department,
        Position:   profile.Position,
        AvatarURL:  profile.AvatarURL,
        Bio:        profile.Bio,
        CreatedAt:  profile.CreatedAt,
        UpdatedAt:  profile.UpdatedAt,
    }
}

func toManagedUserResponse(result *core.ManagedUserResult) ManagedUserResponse {
    return ManagedUserResponse{
        Email:   result.Email,
        Role:    string(result.Role),
        Profile: toUserResponse(&result.Profile),
    }
}

func toListUsersResponse(users []core.UserProfile) ListUsersResponse {
    resp := ListUsersResponse{Users: make([]UserResponse, 0, len(users))}
    for i := range users {
        profile := users[i]
        resp.Users = append(resp.Users, toUserResponse(&profile))
    }
    return resp
}
