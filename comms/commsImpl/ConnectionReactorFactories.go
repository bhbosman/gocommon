package commsImpl

import (
	"context"
	"github.com/bhbosman/gocommon/log"
	"github.com/bhbosman/goerrors"
	"go.uber.org/fx"
)

type ConnectionReactorFactories struct {
	m map[string]IConnectionReactorFactory
}

func (self *ConnectionReactorFactories) CreateClientContext(
	logger *log.SubSystemLogger,
	ConnectionName string,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	name string,
	userContext interface{}) (IConnectionReactor, error) {
	if factory, ok := self.m[name]; ok {
		return factory.Create(ConnectionName, cancelCtx, cancelFunc, logger, userContext), nil
	}
	return nil, goerrors.InvalidParam
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
			return nil, goerrors.DuplicateKey
		}
		m[ccf.Name()] = ccf
	}
	result := &ConnectionReactorFactories{
		m: m,
	}
	return result, nil
}
