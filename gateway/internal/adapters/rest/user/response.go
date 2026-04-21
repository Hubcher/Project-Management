package user

type UserResponse struct {
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

type ManagedUserResponse struct {
    Email   string       `json:"email"`
    Role    string       `json:"role"`
    Profile UserResponse `json:"profile"`
}

type ListUsersResponse struct {
    Users []UserResponse `json:"users"`
}
