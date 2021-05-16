package Bottom

import (
	"github.com/bhbosman/gocomms/internal"
	internalBottom "github.com/bhbosman/gocomms/stacks/Bottom/internal"
	"github.com/google/uuid"
	"github.com/reactivex/rxgo/v2"
)

func StackDefinition(
	connectionType internal.ConnectionType,
	opts ...rxgo.Option) (*internal.StackDefinition, error) {
	id := uuid.New()
	return &internal.StackDefinition{
		IId:      id,
		Name:     internalBottom.StackName,
		Inbound:  internal.NewBoundResultImpl(internalBottom.Inbound(id, opts...)),
		Outbound: internal.NewBoundResultImpl(internalBottom.Outbound(id, opts...)),
	}, nil
}
