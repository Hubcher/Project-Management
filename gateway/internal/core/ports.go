package core

import "context"

type Pinger interface {
	Ping(ctx context.Context) error
}
