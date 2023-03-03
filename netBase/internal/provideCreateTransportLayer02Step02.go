package internal

import (
	"github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func ProvideCreateTransportLayer02Step02() fx.Option {
	return fx.Provide(
		func(
			params struct {
				fx.In
				Logger               *zap.Logger
				CancelCtx            context.Context
				TwoWayPipeDefinition common.ITwoWayPipeDefinition
				StackData            map[string]*common.StackDataContainer
				OutBoundChannel      chan rxgo.Item `name:"OutBoundChannel"`
			},
		) (*common.OutgoingObs, error) {
			params.Logger.Info("createTransportLayer...")
			result, err := params.TwoWayPipeDefinition.BuildOutgoingObs(
				params.OutBoundChannel,
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
