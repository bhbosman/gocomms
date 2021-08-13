package Top

import (
	commsInternal "github.com/bhbosman/gocomms/internal"
	topInternal "github.com/bhbosman/gocomms/stacks/Top/internal"
	"github.com/google/uuid"
	"github.com/reactivex/rxgo/v2"
)

func StackDefinition(
	connectionType commsInternal.ConnectionType,
	opts ...rxgo.Option) (*commsInternal.StackDefinition, error) {
	id := uuid.New()
	return &commsInternal.StackDefinition{
		IId:      id,
		Name:     topInternal.StackName,
		Inbound:  commsInternal.NewBoundResultImpl(topInternal.Inbound(connectionType, id, opts...)),
		Outbound: commsInternal.NewBoundResultImpl(topInternal.Outbound(connectionType, id, opts...)),
	}, nil
}
