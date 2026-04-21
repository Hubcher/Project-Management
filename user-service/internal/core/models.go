package core

import "time"

type User struct {
    ID         string
    FirstName  string
    LastName   string
    MiddleName string
    BirthDate  *time.Time
    Phone      string
    Department string
    Position   string
    AvatarURL  string
    Bio        string
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

type CreateUserInput struct {
    ID         string
    FirstName  string
    LastName   string
    MiddleName string
    BirthDate  *time.Time
    Phone      string
    Department string
    Position   string
    AvatarURL  string
    Bio        string
}

type UpdateUserInput struct {
    ID         string
    FirstName  string
    LastName   string
    MiddleName string
    BirthDate  *time.Time
    Phone      string
    Department string
    Position   string
    AvatarURL  string
    Bio        string
}
