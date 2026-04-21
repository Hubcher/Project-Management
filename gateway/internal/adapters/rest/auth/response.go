package auth

type AuthUserResponse struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
}

type ProfileResponse struct {
    ID         string `json:"id"`
    FirstName  string `json:"first_name"`
    LastName   string `json:"last_name"`
    MiddleName string `json:"middle_name,omitempty"`
    BirthDate  string `json:"birth_date,omitempty"`
    Phone      string `json:"phone,omitempty"`
    Department string `json:"department,omitempty"`
    Position   string `json:"position,omitempty"`
    AvatarURL  string `json:"avatar_url,omitempty"`
    Bio        string `json:"bio,omitempty"`
    CreatedAt  string `json:"created_at,omitempty"`
    UpdatedAt  string `json:"updated_at,omitempty"`
}

type RegisterResponse struct {
    Token   string           `json:"token"`
    User    AuthUserResponse `json:"user"`
    Profile ProfileResponse  `json:"profile"`
}

type LoginResponse struct {
    Token string           `json:"token"`
    User  AuthUserResponse `json:"user"`
}
