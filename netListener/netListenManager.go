package netListener

import (
	"context"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gocomms/impl"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gologging"
	"go.uber.org/fx"
	"golang.org/x/sync/semaphore"
	log2 "log"
	"net"
	"net/url"
)

type netListenManager struct {
	impl.NetManager
	listener interface {
		Accept() (net.Conn, error)
	}
	maxConnections int
}

func (self *netListenManager) listenForNewConnections() {
	go func() {
		n := 0
		sem := semaphore.NewWeighted(int64(self.maxConnections))
		for {
			n++
			self.Logger.LogWithLevel(0, func(logger *log2.Logger) {
				logger.Printf("Trying to accept connections #%v. ", n)
			})
			conn, err := self.Accept()
			if err != nil {
				return
			}
			if sem.TryAcquire(1) {
				self.Logger.LogWithLevel(0, func(logger *log2.Logger) {
					logger.Printf("Accepted connection...")
				})
				conn = impl.NewNetConnWithSemaphoreWrapper(conn, sem)
				self.acceptNewClientConnection(conn)
				continue
			}
			_, _ = conn.Write([]byte("ERR: To many connections\n"))
			_ = conn.Close()
		}
	}()
}

func (self *netListenManager) acceptNewClientConnection(conn net.Conn) {
	go func(conn net.Conn) {
		self.Logger.LogWithLevel(0, func(logger *log2.Logger) {
			logger.Printf("Accepted %s-%s", conn.RemoteAddr(), conn.LocalAddr())
		})
		connectionApp, ctx := self.NewConnectionInstance(internal.ServerConnection, conn)
		err := connectionApp.Err()
		if err != nil {
			_ = conn.Close()
			return
		}
		err = connectionApp.Start(context.Background())
		if err != nil {
			return
		}

		go func(app *fx.App, ctx context.Context) {
			<-ctx.Done()
			app.Stop(context.Background())

		}(connectionApp, ctx)
	}(conn)
}

type NewNetListenAppFunc func(params NetListenAppFuncInParams) (*fx.App, error)

func (self *netListenManager) Accept() (net.Conn, error) {
	return self.listener.Accept()
}

func newNetListenManager(
	params struct {
		fx.In
		Url                        *url.URL
		Listener                   net.Listener
		ConnectionReactorFactories *impl.ConnectionReactorFactories
		ConnectionManager          connectionManager.IConnectionManager
		CancelCtx                  context.Context
		CancelFunction             context.CancelFunc
		StackFactoryFunction       impl.TransportFactoryFunction
		Logger                     *gologging.SubSystemLogger
		ClientContextFactoryName   string `name:"ConnectionReactorFactoryName"`
		LogFactory                 *gologging.Factory
		Cfr                        intf.IConnectionReactorFactory
		Settings                   []ListenAppSettingsApply
	}) *netListenManager {
	netListenSettings := &netListenManagerSettings{
		userContext:    nil,
		maxConnections: 512,
	}
	for _, setting := range params.Settings {
		setting.apply(netListenSettings)
	}
	return &netListenManager{
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
			netListenSettings.userContext),
		listener:       params.Listener,
		maxConnections: netListenSettings.maxConnections,
	}
}
