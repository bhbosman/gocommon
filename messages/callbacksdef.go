package messages

import (
	"context"
	"github.com/bhbosman/goConn"
)

type IApp interface {
	Start(
		//cancellationContext goConn.ICancellationContext,
		ctx context.Context) error
	Stop(
		//cancellationContext goConn.ICancellationContext,
		ctx context.Context) error
	Err() error
}
type CreateAppCallbackFn = func() (IApp, goConn.ICancellationContext, error)

type CreateAppCallback struct {
	Name     string
	Callback CreateAppCallbackFn
}
