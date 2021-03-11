package netDial

import (
	"context"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gologging"
	"go.uber.org/fx"
	"golang.org/x/sync/semaphore"
	"net"
	"net/url"
	"time"
)

type netDialManager struct {
	impl.NetManager
	CanDial       []ICanDial
	MaxConnection int
}

func (self *netDialManager) Start(_ context.Context) error {
	go func() {
		sem := semaphore.NewWeighted(int64(self.MaxConnection))
	loop:
		for {
			if len(self.CanDial) > 0 {
				b := true
				for _, canDial := range self.CanDial {
					b = b && canDial.CanDial()
					if !b {
						time.Sleep(time.Second)
						continue loop
					}
				}
			}

			if self.CancelCtx.Err() != nil {
				return
			}
			if sem.Acquire(self.CancelCtx, 1) != nil {
				return
			}
			dialer := net.Dialer{
				Timeout: time.Second * 30,
			}
			conn, err := dialer.DialContext(self.CancelCtx, "tcp4", self.Url.Host)
			if err != nil {
				sem.Release(1)
				continue loop
			}
			if self.CancelCtx.Err() != nil {
				_ = conn.Close()
				return
			}
			conn = impl.NewNetConnWithSemaphoreWrapper(conn, sem)
			instance, ctx := self.NewConnectionInstance(internal.ClientConnection, conn)
			if instance.Err() != nil {
				_ = conn.Close()
				continue loop
			}
			err = instance.Start(context.Background())
			if err != nil {
				continue loop
			}
			go func(app *fx.App, ctx context.Context) {
				ticker := time.NewTicker(time.Second)
				defer ticker.Stop()
				for {
					select {
					case <-ctx.Done():
						_ = app.Stop(context.Background())
						return
					case <-ticker.C:
						b := true
						for _, canDial := range self.CanDial {
							b = b && canDial.CanDial()
							if !b {
								_ = app.Stop(context.Background())
								return
							}
						}
					}
				}

			}(instance, ctx)
		}
	}()
	return nil
}

func (self *netDialManager) Stop(_ context.Context) error {
	return nil
}

func newNetDialManager(
	params struct {
		fx.In
		Url                        *url.URL
		ConnectionReactorFactories *impl.ConnectionReactorFactories
		ConnectionManager          connectionManager.IConnectionManager
		CancelCtx                  context.Context
		CancelFunction             context.CancelFunc
		StackFactoryFunction       impl.TransportFactoryFunction
		Logger                     *gologging.SubSystemLogger
		ClientContextFactoryName   string `name:"ConnectionReactorFactoryName"`
		LogFactory                 *gologging.Factory
		Options                    []DialAppSettingsApply
		Cfr                        intf.IConnectionReactorFactory
	}) *netDialManager {

	settings := &dialAppSettings{
		userContext: nil,
		canDial:     nil,
	}
	for _, option := range params.Options {
		if option == nil {
			continue
		}
		option.apply(settings)
	}

	return &netDialManager{
		NetManager: impl.NewNetManager(
			params.Url,
			params.ConnectionReactorFactories,
			params.CancelCtx,
			params.CancelFunction,
			params.Logger,
			params.StackFactoryFunction,
			params.ClientContextFactoryName,
			params.Cfr,
			params.ConnectionManager,
			params.LogFactory,
			settings.userContext),
		CanDial:       settings.canDial,
		MaxConnection: settings.maxConnections,
	}
}
