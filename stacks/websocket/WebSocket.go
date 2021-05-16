package websocket

import (
	"context"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/goerrors"
	"github.com/google/uuid"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/multierr"
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

	var data *Data
	data = NewData(stackCancelFunc)
	id := uuid.New()

	// globals
	return &internal.StackDefinition{
		IId:      id,
		Name:     StackName,
		Inbound:  internal.NewBoundResultImpl(inbound(id, data, stackCancelFunc, connectionManager, opts...)),
		Outbound: internal.NewBoundResultImpl(outbound(id, data, stackCancelFunc, connectionManager, opts...)),
		StackState: &internal.StackState{
			Id: id,
			Create: func(Conn net.Conn, Url *url.URL, Ctx context.Context, CancelFunc internal.CancelFunc) (interface{}, error) {
				return internal.NewNoCloser(), nil
			},
			Destroy: func(stackData interface{}) error {
				var err error = nil
				if closer, ok := stackData.(io.Closer); ok {
					err = multierr.Append(err, closer.Close())
				}
				return err
			},
			Start: func(stackData interface{}, startParams internal.StackStartStateParams) (net.Conn, error) {
				return data.OnStart(startParams)
			},
			Stop: func(stackData interface{}, endParams internal.StackEndStateParams) error {
				var err error = nil
				err = multierr.Append(err, data.OnEnd(endParams))
				return err
			},
		},
	}, nil
}
