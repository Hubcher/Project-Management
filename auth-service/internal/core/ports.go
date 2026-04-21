package core

import "context"

type AuthRepository interface {
    CreateAccount(ctx context.Context, acc Account) error
    GetByEmail(ctx context.Context, email string) (Account, error)
    GetByUserID(ctx context.Context, userID string) (Account, error)
    CountAccounts(ctx context.Context) (int, error)
    DeleteByUserID(ctx context.Context, userID string) error
}

type TokenManager interface {
    Generate(claims Claims) (string, error)
    Parse(token string) (Claims, error)
}

type PasswordManager interface {
    Hash(password string) (string, error)
    Compare(hash, password string) error
}
