package internal

import (
	"context"
	"github.com/bhbosman/gocomms/intf"
	"github.com/google/uuid"
	"github.com/reactivex/rxgo/v2"
	"net"
	"net/url"
)

type CancelFunc func(context string, inbound bool, err error)
type PipeDefinitionParams struct {
	ConnectionId      string
	ConnectionManager rxgo.IPublishToConnectionManager
	CancelContext     context.Context
	StackCancelFunc   CancelFunc
	Obs               rxgo.Observable
}

func NewPipeDefinitionParams(
	connectionId string,
	connectionManager rxgo.IPublishToConnectionManager,
	cancelContext context.Context,
	stackCancelFunc CancelFunc,
	obs rxgo.Observable) PipeDefinitionParams {
	return PipeDefinitionParams{
		ConnectionId:      connectionId,
		ConnectionManager: connectionManager,
		CancelContext:     cancelContext,
		StackCancelFunc:   stackCancelFunc,
		Obs:               obs}
}

type PipeDefinition func(params PipeDefinitionParams) (uuid.UUID, rxgo.Observable, error)
type PipeStart func(ctx context.Context) error
type PipeEnd func() error
type PipeState struct {
	Start PipeStart
	End   PipeEnd
}

type StackStartStateParams struct {
	Conn                     net.Conn
	Url                      *url.URL
	Ctx                      context.Context
	CancelFunc               CancelFunc
	ConnectionReactorFactory intf.IConnectionReactorFactoryExtractValues
}

func NewStackStartStateParams(conn net.Conn, url *url.URL, ctx context.Context, cancelFunc CancelFunc, connectionReactorFactoryExtractValues intf.IConnectionReactorFactoryExtractValues) StackStartStateParams {
	return StackStartStateParams{
		Conn:                     conn,
		Url:                      url,
		Ctx:                      ctx,
		CancelFunc:               cancelFunc,
		ConnectionReactorFactory: connectionReactorFactoryExtractValues,
	}
}

type StackStartState func(startParams StackStartStateParams) (net.Conn, error)

type StackEndStateParams struct {
}

func NewStackEndStateParams() StackEndStateParams {
	return StackEndStateParams{}
}

type StackEndState func(endParams StackEndStateParams) error

type StackState struct {
	Start StackStartState
	End   StackEndState
}

type ConnectionType uint8

const (
	ServerConnection ConnectionType = iota
	ClientConnection
)
