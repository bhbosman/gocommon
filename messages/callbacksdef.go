package messages

import (
	"context"
)

type IApp interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Err() error
}
type CreateAppCallbackFn = func() (IApp, context.CancelFunc, error)

type CreateAppCallback struct {
	Name     string
	Callback CreateAppCallbackFn
}
