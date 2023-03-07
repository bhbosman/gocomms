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
					Logger          *zap.Logger
					CancelCtx       context.Context
					PipeDefinition  common.IPipeDefinition `name:"Outbound"`
					StackData       map[string]*common.StackDataContainer
					OutBoundChannel *rxgo.ItemChannel `name:"OutBoundChannel"`
				},
			) (gocommon.IObservable, error) {
				params.Logger.Info("createTransportLayer...")
				result, err := params.PipeDefinition.BuildObs(
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
