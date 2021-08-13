package messageCompressor

import (
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/goerrors"
	"github.com/google/uuid"
	"github.com/reactivex/rxgo/v2"
)

func StackDefinition(
	connectionType internal.ConnectionType,
	stackCancelFunc internal.CancelFunc,
	opts ...rxgo.Option) (*internal.StackDefinition, error) {
	if stackCancelFunc == nil {
		return nil, goerrors.InvalidParam
	}

	id := uuid.New()
	return &internal.StackDefinition{
		IId:      id,
		Name:     StackName,
		Inbound:  internal.NewBoundResultImpl(inbound(id, stackCancelFunc, opts...)),
		Outbound: internal.NewBoundResultImpl(outbound(id, stackCancelFunc, opts...)),
	}, nil
}
