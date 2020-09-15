package commsImpl

import (
	"context"
	"github.com/bhbosman/gocommon/constants"
	"go.uber.org/fx"
)

type ConnectionReactorFactories struct {
	m map[string]IConnectionReactorFactory
}

func (self *ConnectionReactorFactories) CreateClientContext(
	logger fx.ILogger,
	ConnectionName string,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	name string,
	userContext interface{}) (IConnectionReactor, error) {
	if factory, ok := self.m[name]; ok {
		return factory.Create(ConnectionName, cancelCtx, cancelFunc, logger, userContext), nil
	}
	return nil, constants.InvalidParam
}

const ConnectionReactorFactoryConst = "ConnectionReactorFactory"

func newConnectionReactorFactories(
	params struct {
		fx.In
		Factories []IConnectionReactorFactory `group:"ConnectionReactorFactory"`
	}) (*ConnectionReactorFactories, error) {
	m := make(map[string]IConnectionReactorFactory)
	for _, ccf := range params.Factories {
		if _, ok := m[ccf.Name()]; ok {
			return nil, constants.DuplicateKey
		}
		m[ccf.Name()] = ccf
	}
	result := &ConnectionReactorFactories{
		m: m,
	}
	return result, nil
}
