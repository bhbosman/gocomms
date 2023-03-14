package netBase

import (
	"context"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocomms/common"
	"go.uber.org/fx"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"net/url"
	"sync"
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

func (self *NetManager) RegisterConnectionShutdown(
	connectionId string,
	callback func(),
	cancellationContext ...common.ICancellationContext,
) error {
	mutex := sync.Mutex{}
	cancelCalled := false
	cb := func() {
		mutex.Lock()
		b := cancelCalled
		cancelCalled = true
		mutex.Unlock()
		if !b {
			callback()
		}
		for _, instance := range cancellationContext {
			_ = instance.Remove(connectionId)
		}
	}
	var result error
	for _, ctx := range cancellationContext {
		b, err := ctx.Add(connectionId, cb)
		result = multierr.Append(result, err)
		if !b {
			cb()
			return result
		}
	}
	return result
}
