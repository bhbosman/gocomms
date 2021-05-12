package Top

import (
	"github.com/bhbosman/gocomms/internal"
	internal2 "github.com/bhbosman/gocomms/stacks/Top/internal"
	"github.com/google/uuid"
	"github.com/reactivex/rxgo/v2"
)

func StackDefinition(
	connectionType internal.ConnectionType,
	opts ...rxgo.Option) (*internal.StackDefinition, error) {
	id := uuid.New()
	return &internal.StackDefinition{
		Id:       id,
		Name:     internal2.StackName,
		Inbound:  internal.NewBoundResultImpl(internal2.Inbound(id, opts...)),
		Outbound: internal.NewBoundResultImpl(internal2.Outbound(id, opts...)),
	}, nil
}
