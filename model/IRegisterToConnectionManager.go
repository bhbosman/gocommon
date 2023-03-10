package model

import (
	"context"
	"github.com/reactivex/rxgo/v2"
)

type IRegisterToConnectionManager interface {
	RegisterConnection(
		id string,
		function context.CancelFunc,
		CancelContext context.Context,
		nextFuncOutBoundChannel rxgo.NextFunc,
		nextFuncInBoundChannel rxgo.NextFunc,
	) error
	DeregisterConnection(id string) error
	NameConnection(id string, name string) error
}
