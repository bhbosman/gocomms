package Bottom

import (
	"github.com/bhbosman/gocomms/internal"
	internalBottom "github.com/bhbosman/gocomms/stacks/Bottom/internal"
	"github.com/google/uuid"
	"github.com/reactivex/rxgo/v2"
	"net"
)

func StackDefinition(
	connectionType internal.ConnectionType,
	opts ...rxgo.Option) (*internal.StackDefinition, error) {
	id := uuid.New()
	return &internal.StackDefinition{
		Id:       id,
		Name:     internalBottom.StackName,
		Inbound:  internal.NewBoundResultImpl(internalBottom.Inbound(id, opts...)),
		Outbound: internal.NewBoundResultImpl(internalBottom.Outbound(id, opts...)),
		StackState: internal.StackState{
			Start: func(startParams internal.StackStartStateParams) (net.Conn, error) {
				return startParams.Conn, startParams.Ctx.Err()
			},
		},
	}, nil
}
