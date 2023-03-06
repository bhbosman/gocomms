package internal

import (
	"github.com/bhbosman/gocommon"
	"github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func ProvideCreateTransportLayer02Step02() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Name: "Outbound",
			Target: func(
				params struct {
					fx.In
					Logger                 *zap.Logger
					CancelCtx              context.Context
					OutboundPipeDefinition common.IOutboundPipeDefinition
					StackData              map[string]*common.StackDataContainer
					OutBoundChannel        chan rxgo.Item `name:"OutBoundChannel"`
				},
			) (gocommon.IObservable, error) {
				params.Logger.Info("createTransportLayer...")
				result, err := params.OutboundPipeDefinition.BuildOutgoingObs(
					params.OutBoundChannel,
					params.StackData,
					params.CancelCtx,
				)
				if err != nil {
					return nil, err
				}
				return result, nil
			},
		},
	)
}
