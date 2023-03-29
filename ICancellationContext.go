package gocommon

import (
	"context"
)

type ICancellationContext interface {
	Add(connectionId string, f func()) (bool, error)
	Remove(connectionId string) error
	Cancel(string)
	CancelWithError(string, error)
	Err() error
	CancelContext() context.Context
	CancelFunc() context.CancelFunc
}

type IApp interface {
	Start(
		ctx context.Context) error
	Stop(
		ctx context.Context) error
	Err() error
}
type CreateAppCallbackFn = func() (IApp, ICancellationContext, error)

type CreateAppCallback struct {
	Name     string
	Callback CreateAppCallbackFn
}
