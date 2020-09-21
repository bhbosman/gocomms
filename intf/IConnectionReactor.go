package intf

import (
	"context"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gologging"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
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
	Create(name string, cancelCtx context.Context, cancelFunc context.CancelFunc, logger *gologging.SubSystemLogger, userContext interface{}) IConnectionReactor
}
