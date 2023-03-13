package netBase

import (
	"context"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocomms/common"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/url"
)

type NetManager struct {
	CancellationContext                      common.ICancellationContext
	CancelCtx                                context.Context
	ZapLogger                                *zap.Logger
	ConnectionManager                        goConnectionManager.IService
	ConnectionUrl                            *url.URL
	UseProxy                                 bool
	ProxyUrl                                 *url.URL
	Name                                     string
	ConnectionInstancePrefix                 string
	UniqueSessionNumber                      interfaces.IUniqueReferenceService
	AdditionalFxOptionsForConnectionInstance func() fx.Option
	GoFunctionCounter                        GoFunctionCounter.IService
}

func NewNetManager(
	name string,
	connectionInstancePrefix string,
	useProxy bool,
	proxyUrl *url.URL,
	connectionUrl *url.URL,
	cancelCtx context.Context,
	CancellationContext common.ICancellationContext,
	connectionManager goConnectionManager.IService,
	ZapLogger *zap.Logger,
	uniqueSessionNumber interfaces.IUniqueReferenceService,
	additionalFxOptionsForConnectionInstance func() fx.Option,
	GoFunctionCounter GoFunctionCounter.IService,
) (NetManager, error) {
	return NetManager{
		CancellationContext:                      CancellationContext,
		CancelCtx:                                cancelCtx,
		ZapLogger:                                ZapLogger,
		ConnectionManager:                        connectionManager,
		ConnectionUrl:                            connectionUrl,
		UseProxy:                                 useProxy,
		ProxyUrl:                                 proxyUrl,
		Name:                                     name,
		ConnectionInstancePrefix:                 connectionInstancePrefix,
		UniqueSessionNumber:                      uniqueSessionNumber,
		AdditionalFxOptionsForConnectionInstance: additionalFxOptionsForConnectionInstance,
		GoFunctionCounter:                        GoFunctionCounter,
	}, nil
}
