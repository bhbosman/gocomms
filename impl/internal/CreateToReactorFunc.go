package internal

import (
	"context"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/goprotoextra"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
	"time"
)

func CreateToReactorFunc(
	params struct {
		fx.In
		CancelCtx context.Context
		Ch        *internal.ChannelManager
	}) goprotoextra.ToReactorFunc {
	return func(inline bool, any interface{}) error {
		err := params.CancelCtx.Err()
		if err != nil {
			return err
		}
		if params.Ch.Active() {
			if !inline {
				go func(any interface{}) {
					params.Ch.Send(params.CancelCtx, rxgo.NewNextExternal(false, any))
				}(any)
				return nil
			}
			return params.Ch.SendContextWithTimeOutAndRetries(
				rxgo.NewNextExternal(false, any),
				params.CancelCtx,
				time.Millisecond*500,
				5,
				params.Ch.Items)
		}
		return params.CancelCtx.Err()
	}
}
