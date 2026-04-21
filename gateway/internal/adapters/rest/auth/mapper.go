package auth

import "github.com/Hubcher/project-management/gateway/internal/core"

func toRegisterInput(req RegisterRequest) core.RegisterInput {
    return core.RegisterInput{
        Email:      req.Email,
        Password:   req.Password,
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

func toAuthUserResponse(user core.AuthUser) AuthUserResponse {
    return AuthUserResponse{
        UserID: user.UserID,
        Email:  user.Email,
        Role:   string(user.Role),
    }
}

func toProfileResponse(profile core.UserProfile) ProfileResponse {
    return ProfileResponse{
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

func toRegisterResponse(result *core.RegisterResult) RegisterResponse {
    return RegisterResponse{
        Token:   result.Token,
        User:    toAuthUserResponse(result.User),
        Profile: toProfileResponse(result.Profile),
    }
}

func toLoginResponse(result *core.LoginResult) LoginResponse {
    return LoginResponse{
        Token: result.Token,
        User:  toAuthUserResponse(result.User),
    }
}
