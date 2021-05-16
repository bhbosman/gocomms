package internal

import (
	"context"
	"github.com/bhbosman/gocomms/internal"
	"github.com/reactivex/rxgo/v2"
	"go.uber.org/fx"
)

func CreateChannel(
	params struct {
		fx.In
		CancelCtx    context.Context
		ConnectionId string `name:"ConnectionId"`
		Lifecycle    fx.Lifecycle
	}) (rxgo.Observable, *internal.ChannelManager, error) {
	ch := internal.NewChannelManager("create channel", params.ConnectionId)
	obs := rxgo.FromChannel(ch.Items, rxgo.WithContext(params.CancelCtx))
	params.Lifecycle.Append(fx.Hook{
		OnStart: nil,
		OnStop: func(ctx context.Context) error {
			return ch.Close()
		},
	})
	return obs, ch, nil
}
