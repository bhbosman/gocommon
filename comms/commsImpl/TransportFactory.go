package commsImpl

import (
	"context"
	"github.com/bhbosman/gocommon/comms/common"
	"github.com/bhbosman/gocommon/comms/connectionManager"
	"github.com/bhbosman/gocommon/stacks/Bottom"
	"github.com/bhbosman/gocommon/stacks/KillConnection"
	"github.com/bhbosman/gocommon/stacks/Top"
	"github.com/bhbosman/gocommon/stacks/defs"
	"github.com/bhbosman/gocommon/stacks/messageBreaker"
	"github.com/bhbosman/gocommon/stacks/messageCompressor"
	"github.com/bhbosman/gocommon/stacks/messageNumber"
	"github.com/bhbosman/gocommon/stacks/tlsConnection"
	"github.com/bhbosman/gocommon/stacks/websocket"
	"github.com/reactivex/rxgo/v2"
)

type TransportFactoryFunction func(
	connectionType common.ConnectionType,
	connectionId string,
	userContext interface{},
	connectionManager connectionManager.IConnectionManager,
	cancelContext context.Context,
	cancelFunc defs.CancelFunc,
	opts ...rxgo.Option) (*defs.TwoWayPipeDefinition, error)

type TransportFactory struct {
}

func (self *TransportFactory) CreateEmpty(
	connectionType common.ConnectionType,
	_ string,
	_ interface{},
	_ connectionManager.IConnectionManager,
	_ context.Context,
	_ defs.CancelFunc,
	opts ...rxgo.Option) (*defs.TwoWayPipeDefinition, error) {
	result := defs.NewTwoWayPipeDefinition(nil)

	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return Top.StackDefinition()
	})
	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return Bottom.StackDefinition()
	})
	return result, nil
}

func (self *TransportFactory) CreateWebSocket(
	connectionType common.ConnectionType,
	connectionId string,
	userContext interface{},
	connectionManager connectionManager.IConnectionManager,
	cancelContext context.Context,
	stackCancelFunc defs.CancelFunc,
	opts ...rxgo.Option) (*defs.TwoWayPipeDefinition, error) {
	result := defs.NewTwoWayPipeDefinition(nil)

	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return Top.StackDefinition()
	})
	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return websocket.StackDefinition(stackCancelFunc, connectionManager, opts...)
	})
	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return Bottom.StackDefinition()
	})
	return result, nil
}

func (self *TransportFactory) CreateDefault(
	connectionType common.ConnectionType,
	_ string,
	userContext interface{},
	connectionManager connectionManager.IConnectionManager,
	cancelContext context.Context,
	stackCancelFunc defs.CancelFunc,
	opts ...rxgo.Option) (*defs.TwoWayPipeDefinition, error) {
	result := defs.NewTwoWayPipeDefinition(nil)

	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return Top.StackDefinition()
	})

	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return KillConnection.StackDefinition(cancelContext, stackCancelFunc, connectionManager, opts...)
	})
	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return Bottom.StackDefinition()
	})
	return result, nil
}

func (self *TransportFactory) CreateCompressed(
	connectionType common.ConnectionType,
	connectionId string,
	userContext interface{},
	connectionManager connectionManager.IConnectionManager,
	_ context.Context,
	stackCancelFunc defs.CancelFunc,
	opts ...rxgo.Option) (*defs.TwoWayPipeDefinition, error) {
	result := defs.NewTwoWayPipeDefinition(nil)
	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return Top.StackDefinition()
	})
	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return messageCompressor.StackDefinition(stackCancelFunc, opts...)
	})
	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return messageNumber.StackDefinition(userContext, stackCancelFunc, opts...)
	})
	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return messageBreaker.StackDefinition(stackCancelFunc, nil, connectionManager, opts...)
	})
	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return Bottom.StackDefinition()
	})
	return result, nil
}

func (self *TransportFactory) CreateUnCompressed(
	connectionType common.ConnectionType,
	_ string,
	userContext interface{},
	connectionManager connectionManager.IConnectionManager,
	_ context.Context,
	stackCancelFunc defs.CancelFunc,
	opts ...rxgo.Option) (*defs.TwoWayPipeDefinition, error) {
	result := defs.NewTwoWayPipeDefinition(nil)
	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return Top.StackDefinition()
	})
	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return messageNumber.StackDefinition(userContext, stackCancelFunc, opts...)
	})
	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return messageBreaker.StackDefinition(stackCancelFunc, nil, connectionManager, opts...)
	})
	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return Bottom.StackDefinition()
	})
	return result, nil
}

func (self *TransportFactory) CreateCompressedTls(
	connectionType common.ConnectionType,
	connectionId string,
	userContext interface{},
	connectionManager connectionManager.IConnectionManager,
	_ context.Context,
	stackCancelFunc defs.CancelFunc,
	opts ...rxgo.Option) (*defs.TwoWayPipeDefinition, error) {
	result := defs.NewTwoWayPipeDefinition(nil)
	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return Top.StackDefinition()
	})
	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return messageCompressor.StackDefinition(stackCancelFunc, opts...)
	})
	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return messageNumber.StackDefinition(userContext, stackCancelFunc, opts...)
	})
	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return messageBreaker.StackDefinition(stackCancelFunc, nil, connectionManager, opts...)
	})
	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return tlsConnection.StackDefinition(connectionType, stackCancelFunc, connectionManager, connectionId)
	})
	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return Bottom.StackDefinition()
	})
	return result, nil
}

func (self *TransportFactory) CreateUnCompressedTls(
	connectionType common.ConnectionType,
	connectionId string,
	userContext interface{},
	connectionManager connectionManager.IConnectionManager,
	_ context.Context,
	stackCancelFunc defs.CancelFunc,
	opts ...rxgo.Option) (*defs.TwoWayPipeDefinition, error) {
	result := defs.NewTwoWayPipeDefinition(nil)
	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return Top.StackDefinition()
	})
	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return messageNumber.StackDefinition(userContext, stackCancelFunc, opts...)
	})
	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return messageBreaker.StackDefinition(stackCancelFunc, nil, connectionManager, opts...)
	})
	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return tlsConnection.StackDefinition(connectionType, stackCancelFunc, connectionManager, connectionId)
	})
	result.AddStackDefinitionFunc(func() (*defs.StackDefinition, error) {
		return Bottom.StackDefinition()
	})
	return result, nil
}

const TransportFactoryCompressedTlsName = "CompressedTLS"
const TransportFactoryUnCompressedTlsName = "UncompressedTLS"
const TransportFactoryCompressedName = "Compressed"
const TransportFactoryUnCompressedName = "Uncompressed"
const TransportFactoryEmptyName = "Empty"
const WebSocketName = "WebSocket"

func (self *TransportFactory) Get(name string) (TransportFactoryFunction, error) {
	switch name {
	case TransportFactoryCompressedTlsName:
		return self.CreateCompressedTls, nil
	case TransportFactoryUnCompressedTlsName:
		return self.CreateUnCompressedTls, nil
	case TransportFactoryCompressedName:
		return self.CreateCompressed, nil
	case TransportFactoryUnCompressedName:
		return self.CreateUnCompressed, nil
	case TransportFactoryEmptyName:
		return self.CreateEmpty, nil
	case WebSocketName:
		return self.CreateWebSocket, nil
	default:
		return self.CreateDefault, nil
	}
}

func NewTransportFactory() *TransportFactory {
	result := &TransportFactory{}
	return result
}
