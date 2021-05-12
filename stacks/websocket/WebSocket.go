package websocket

import (
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/goerrors"
	"github.com/google/uuid"
	"github.com/reactivex/rxgo/v2"
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
		Id:       id,
		Name:     StackName,
		Inbound:  internal.NewBoundResultImpl(inbound(id, data, stackCancelFunc, connectionManager, opts...)),
		Outbound: internal.NewBoundResultImpl(outbound(id, data, stackCancelFunc, connectionManager, opts...)),
		StackState: internal.StackState{
			Start: data.OnStart,
			End:   data.OnEnd,
		},
	}, nil
}
