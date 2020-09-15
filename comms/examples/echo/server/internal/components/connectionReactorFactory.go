package Server

import (
	"context"
	"github.com/bhbosman/gocommon/comms/commsImpl"
	"go.uber.org/fx"
)

type connectionReactorFactory struct {
	name string
}

func (self *connectionReactorFactory) Create(
	name string,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	logger fx.ILogger,
	userContext interface{}) commsImpl.IConnectionReactor {
	result := newConnectionReactor(logger, cancelCtx, cancelFunc, name, userContext)
	return result
}

func (self *connectionReactorFactory) Name() string {
	return self.name
}
