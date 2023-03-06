package internal

import (
	"github.com/bhbosman/gocommon"
	"github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func ProvideCreateTransportLayer02Step01() fx.Option {
	return fx.Provide(
		fx.Annotated{
			Name: "Inbound",
			Target: func(
				params struct {
					fx.In
					Logger                *zap.Logger
					CancelCtx             context.Context
					InboundPipeDefinition common.IPipeDefinition `name:"Inbound"`
					StackData             map[string]*common.StackDataContainer
					InBoundChannel        chan rxgo.Item `name:"InBoundChannel"`
				},
			) (gocommon.IObservable, error) {
				params.Logger.Info("createTransportLayer...")

				result, err := params.InboundPipeDefinition.BuildOutgoingObs(
					params.InBoundChannel,
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
