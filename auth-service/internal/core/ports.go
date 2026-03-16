package core

import "context"

type DB interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid string, err error)
	User(ctx context.Context, email string) (User, error)
	isAdmin(ctx context.Context, userID string) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (App, error)
}
