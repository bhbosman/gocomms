package internal

import (
	"fmt"
	"github.com/bhbosman/gocommon"
	"github.com/bhbosman/gocommon/model"
	"github.com/bhbosman/gocomms/RxHandlers"
	"github.com/bhbosman/gocomms/common"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func ProvideChannel(name string) fx.Option {
	return fx.Provide(
		fx.Annotate(
			func(
				_ fx.Lifecycle,
				pipeDefinition common.IPipeDefinition,
				stackData map[string]*common.StackDataContainer,
				cancelCtx context.Context,
				connectionId string,
				cancellationContext gocommon.ICancellationContext,
			) (chan rxgo.Item, gocommon.IObservable, error) {
				ch := make(chan rxgo.Item)
				result, err := pipeDefinition.BuildObs(ch, stackData, cancelCtx)
				cancellationContext.Add(
					fmt.Sprintf("%v_%v", connectionId, name),
					func(ch chan rxgo.Item) func() {
						return func() {
							close(ch)
						}
					}(ch))
				if err != nil {
					close(ch)
					return nil, nil, err
				}
				return ch, result, nil
			},
			fx.ParamTags(
				``,
				fmt.Sprintf(`name:"%v"`, name),
				``,
				``,
				`name:"ConnectionId"`,
				``,
			),
			fx.ResultTags(
				fmt.Sprintf(`name:"%v"`, name),
				fmt.Sprintf(`name:"%v"`, name),
			),
		),
		fx.Annotate(
			func(
				ch chan rxgo.Item,
				logger *zap.Logger,
				ctx context.Context,
				connectionId string,
			) (rxgo.NextFunc, rxgo.ErrFunc, rxgo.CompletedFunc, error) {
				handler, err := RxHandlers.All2(
					fmt.Sprintf("ProvideChannel %v", connectionId),
					model.StreamDirectionUnknown,
					ch,
					logger,
					ctx,
					true,
				)
				if err != nil {
					return nil, nil, nil, err
				}
				return handler.OnSendData, handler.OnError, handler.OnComplete, nil
			},
			fx.ParamTags(
				fmt.Sprintf(`name:"%v"`, name),
				``,
				``,
				`name:"ConnectionId"`,
			),
			fx.ResultTags(
				fmt.Sprintf(`name:"%v"`, name),
				fmt.Sprintf(`name:"%v"`, name),
				fmt.Sprintf(`name:"%v"`, name),
			),
		),
	)
}
