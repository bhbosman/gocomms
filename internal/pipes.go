package internal

import (
	"context"
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

type PipeDefinition func(params PipeDefinitionParams) (rxgo.Observable, error)
type PipeStart func(ctx context.Context) error
type PipeEnd func() error
type PipeState struct {
	Start PipeStart
	End   PipeEnd
}

type StackStartState func(conn net.Conn, url *url.URL, ctx context.Context, cancelFunc CancelFunc) (net.Conn, error)
type StackEndState func() error
type StackState struct {
	Start StackStartState
	End   StackEndState
}

type StackDefinition struct {
	Name       string
	Inbound    func(index int, ctx context.Context) BoundDefinition
	Outbound   func(index int, ctx context.Context) BoundDefinition
	StackState StackState
}




type ConnectionType uint8

const (
	ServerConnection ConnectionType = iota
	ClientConnection
)
