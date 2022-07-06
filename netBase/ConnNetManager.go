package netBase

import (
	"context"
	"fmt"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/interfaces"
	fx2 "github.com/bhbosman/gocommon/fx"
	"github.com/bhbosman/gocommon/model"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"net"
	"net/url"
	"time"
)

type ConnNetManager struct {
	NetManager
}

func NewConnNetManager(
	name string,
	connectionInstancePrefix string,
	useProxy bool,
	proxyUrl *url.URL,
	connectionUrl *url.URL,
	cancelCtx context.Context,
	cancelFunction context.CancelFunc,
	//stackName string,
	connectionManager goConnectionManager.IService,
	userContext interface{},
	ZapLogger *zap.Logger,
	uniqueSessionNumber interfaces.IUniqueReferenceService,
	additionalFxOptionsForConnectionInstance func() fx.Option,
	GoFunctionCounter GoFunctionCounter.IService,
) (ConnNetManager, error) {

	netManager, err := NewNetManager(
		name,
		connectionInstancePrefix,
		useProxy,
		proxyUrl,
		connectionUrl,
		cancelCtx,
		cancelFunction,
		//stackName,
		connectionManager,
		userContext,
		ZapLogger,
		uniqueSessionNumber,
		additionalFxOptionsForConnectionInstance,
		GoFunctionCounter,
	)
	if err != nil {
		return ConnNetManager{}, err
	}
	return ConnNetManager{
		NetManager: netManager,
	}, nil
}

// NewConnectionInstance is the default behaviour, which will use the StackName the connection was assigned with.
func (self *ConnNetManager) NewConnectionInstance(
	uniqueReference string,
	goFunctionCounter GoFunctionCounter.IService,
	connectionType model.ConnectionType,
	conn net.Conn,
) (*fx.App, context.Context, context.CancelFunc) {
	return self.NewConnectionInstanceWithStackName(
		uniqueReference,
		goFunctionCounter,
		connectionType,
		conn,
		//self.StackName,
	)
}

// NewConnectionInstanceWithStackName give the ability to override the stack name, after it has been set
// This is useful in SSH, when ssh can dynamically assign a stack name, depending on the ssh protocol required
func (self *ConnNetManager) NewConnectionInstanceWithStackName(
	uniqueReference string,
	goFunctionCounter GoFunctionCounter.IService,
	connectionType model.ConnectionType,
	conn net.Conn,
	//stackName string,
	settingOptions ...INewConnectionInstanceSettingsApply,
) (*fx.App, context.Context, context.CancelFunc) {
	var resultContext context.Context
	var resultCancelFunc context.CancelFunc
	fxAppOptions := self.NewConnectionInstanceOptions(
		uniqueReference,
		goFunctionCounter,
		connectionType,
		conn,
		//stackName,
		NewAddFxOptions(fx.Populate(&resultContext)),
		NewAddFxOptions(fx.Populate(&resultCancelFunc)),
		newAddSettings(settingOptions...),
	)
	fxApp := fx.New(fxAppOptions)
	err := fxApp.Err()
	if err != nil {
		if resultCancelFunc != nil {
			resultCancelFunc()
		}
		if conn != nil {
			// this is required, when providers fail, and the Conn.Close() invoker has not been added to the
			//fx.App lifecycle stack for destruction
			// this error can be ignored
			err = multierr.Append(err, conn.Close())
		}
		if err != nil {
			self.ZapLogger.Error(
				"On a presumed error in the providers, the connection may double close",
				zap.Error(err))
		}

	}
	return fxApp, resultContext, resultCancelFunc
}

func (self *ConnNetManager) NewConnectionInstanceOptions(
	uniqueReference string,
	goFunctionCounter GoFunctionCounter.IService,
	connectionType model.ConnectionType,
	conn net.Conn,
	//stackName string,
	settingOptions ...INewConnectionInstanceSettingsApply,
) fx.Option {
	rwcOptions := self.NewReaderWriterCloserInstanceOptions(
		uniqueReference,
		goFunctionCounter,
		connectionType,
		conn,
		//stackName,
		fx.Provide(
			func() (*zap.Logger, error) {
				if self.ZapLogger == nil {
					return nil, fmt.Errorf("zap.Logger is nil. Please resolve")
				}
				return self.ZapLogger.With(
						zap.String("UniqueReference", uniqueReference),
						zap.String("LocalAddr", conn.LocalAddr().String()),
						zap.String("RemoteAddr", conn.RemoteAddr().String())),
					nil
			},
		),
		settingOptions...,
	)
	return fx2.NewFxApplicationOptions(
		time.Hour,
		time.Hour,
		rwcOptions,
		fx.Provide(
			fx.Annotated{
				Target: func() (net.Conn, error) {
					if conn == nil {
						return nil, fmt.Errorf("connection is nil. Please resolve")
					}
					return conn, nil
				},
			},
		),
	)
}
