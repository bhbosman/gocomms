package impl

import (
	"context"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gologging"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
	"net"
	"net/url"
)

type BaseConnectionReactor struct {
	// Name is the name of the connection.
	Name string
	// CancelCtx is the cancellation context associated with the connection. This can be used to check if the connection
	// have not been closed
	CancelCtx         context.Context
	CancelFunc        context.CancelFunc
	Logger            *gologging.SubSystemLogger
	ToConnection      goprotoextra.ToConnectionFunc
	ToReactor         goprotoextra.ToReactorFunc
	Conn              net.Conn
	ConnectionId      string
	Url               *url.URL
	UrlAsString       string
	ConnectionManager connectionManager.IConnectionManager
	UserContext       interface{}
}

func NewBaseConnectionReactor(
	logger *gologging.SubSystemLogger,
	name string,
	cancelCtx context.Context,
	cancelFunc context.CancelFunc,
	userContext interface{}) BaseConnectionReactor {
	return BaseConnectionReactor{
		Name:         name,
		CancelCtx:    cancelCtx,
		CancelFunc:   cancelFunc,
		Logger:       logger,
		ToConnection: nil,
		UserContext:  userContext,
	}
}

func (self *BaseConnectionReactor) Init(
	conn net.Conn,
	url *url.URL,
	connectionId string,
	connectionManager connectionManager.IConnectionManager,
	toConnectionFunc goprotoextra.ToConnectionFunc,
	toConnectionReactor goprotoextra.ToReactorFunc) (rxgo.NextExternalFunc, error) {
	self.Conn = conn
	self.Url = url
	self.UrlAsString = url.String()
	self.ConnectionManager = connectionManager
	self.ConnectionId = connectionId
	self.ToReactor = toConnectionReactor
	self.ToConnection = toConnectionFunc

	return func(external bool, i interface{}) {

	}, nil
}

func (self *BaseConnectionReactor) Close() error {
	return nil
}

func (self *BaseConnectionReactor) Open() error {
	return nil
}

func (self *BaseConnectionReactor) SendStringToConnection(s string) error {
	rws, err := gomessageblock.NewReaderWriterString(s)
	if err != nil {
		return err
	}
	return self.ToConnection(rws)
}
