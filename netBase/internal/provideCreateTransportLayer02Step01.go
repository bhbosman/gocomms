package internal

import (
	"github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func ProvideCreateTransportLayer02Step01() fx.Option {
	return fx.Provide(
		func(
			params struct {
				fx.In
				Logger               *zap.Logger
				CancelCtx            context.Context
				TwoWayPipeDefinition *common.TwoWayPipeDefinition
				StackData            map[string]*common.StackDataContainer
				InBoundChannel       chan rxgo.Item `name:"InBoundChannel"`
			},
		) (*common.IncomingObs, error) {
			params.Logger.Info("createTransportLayer...")

			result, err := params.TwoWayPipeDefinition.BuildIncomingObs(
				params.InBoundChannel,
				params.StackData,
				params.CancelCtx,
			)
			if err != nil {
				return nil, err
			}
			return result, nil
		},
	)
}
