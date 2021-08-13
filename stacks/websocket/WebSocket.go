package websocket

import (
	"context"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/goerrors"
	"github.com/google/uuid"
	"github.com/reactivex/rxgo/v2"
	"io"
	"net"
	"net/url"
)

func StackDefinition(
	connectionType internal.ConnectionType,
	stackCancelFunc internal.CancelFunc,
	connectionManager rxgo.IPublishToConnectionManager,
	opts ...rxgo.Option) (*internal.StackDefinition, error) {
	if stackCancelFunc == nil {
		return nil, goerrors.InvalidParam
	}
	id := uuid.New()
	return &internal.StackDefinition{
		IId:      id,
		Name:     StackName,
		Inbound:  internal.NewBoundResultImpl(inbound(id, connectionType, stackCancelFunc, connectionManager, opts...)),
		Outbound: internal.NewBoundResultImpl(outbound(id, connectionType, stackCancelFunc, connectionManager, opts...)),
		StackState: &internal.StackState{
			Id: id,
			Create: func(Conn net.Conn, Url *url.URL, Ctx context.Context, CancelFunc internal.CancelFunc, cfr intf.IConnectionReactorFactoryExtractValues) (interface{}, error) {
				return NewStackData(Conn, stackCancelFunc, Ctx, Url, cfr)
			},
			Destroy: func(stackData interface{}) error {
				if _, ok := stackData.(*StackData); !ok {
					return WrongStackDataError(connectionType, stackData)
				}
				if closer, ok := stackData.(io.Closer); ok {
					return closer.Close()
				}
				return nil
			},
			Start: func(stackData interface{}, startParams internal.StackStartStateParams) (net.Conn, error) {
				sd, ok := stackData.(*StackData)
				if !ok {
					return nil, WrongStackDataError(connectionType, stackData)
				}
				return sd.OnStart(startParams)
			},
			Stop: func(stackData interface{}, endParams internal.StackEndStateParams) error {
				_, ok := stackData.(*StackData)
				if !ok {
					return WrongStackDataError(connectionType, stackData)
				}
				var err error = nil
				return err
			},
		},
	}, nil
}
