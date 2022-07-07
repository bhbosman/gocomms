package internal

import (
	"github.com/bhbosman/gocomms/common"
	"go.uber.org/fx"
)

func ProvideCreateTransportLayer00() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func(
				params struct {
					fx.In
					TwoWayPipeDefinition *common.TwoWayPipeDefinition
				},
			) ([]*common.StackState, error) {
				return params.TwoWayPipeDefinition.BuildStackState()
			},
		},
	)
}
