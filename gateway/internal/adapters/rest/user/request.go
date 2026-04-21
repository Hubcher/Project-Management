package user

type CreateUserRequest struct {
    Email      string `json:"email"`
    Password   string `json:"password"`
    Role       string `json:"role"`
    FirstName  string `json:"first_name"`
    LastName   string `json:"last_name"`
    MiddleName string `json:"middle_name"`
    BirthDate  string `json:"birth_date"`
    Phone      string `json:"phone"`
    Department string `json:"department"`
    Position   string `json:"position"`
    AvatarURL  string `json:"avatar_url"`
    Bio        string `json:"bio"`
}

type UpdateUserRequest struct {
    FirstName  string `json:"first_name"`
    LastName   string `json:"last_name"`
    MiddleName string `json:"middle_name"`
    BirthDate  string `json:"birth_date"`
    Phone      string `json:"phone"`
    Department string `json:"department"`
    Position   string `json:"position"`
    AvatarURL  string `json:"avatar_url"`
    Bio        string `json:"bio"`
}
