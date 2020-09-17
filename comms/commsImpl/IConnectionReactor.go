package commsImpl

import (
	"context"
	rxgo "github.com/ReactiveX/RxGo"
	"github.com/bhbosman/gocommon/comms/connectionManager"
	"github.com/bhbosman/gocommon/log"
	"github.com/bhbosman/goprotoextra"
	"net"
	"net/url"
)

type IConnectionReactor interface {
	Init(
		conn net.Conn,
		url *url.URL,
		connectionId string,
		connectionManager connectionManager.IConnectionManager,
		onSend goprotoextra.ToConnectionFunc,
		toConnectionReactor goprotoextra.ToReactorFunc) (rxgo.NextExternalFunc, error)
	Close() error
	Open() error
}

const ConnectionName = "ConnectionName"
const ConnectionId = "ConnectionId"
const UserContext = "UserContext"
const ConnectionReactorFactoryName = "ConnectionReactorFactoryName"

type IConnectionReactorFactory interface {
	Name() string
	Create(name string, cancelCtx context.Context, cancelFunc context.CancelFunc, logger *log.SubSystemLogger, userContext interface{}) IConnectionReactor
}
