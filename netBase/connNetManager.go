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
	CancellationContext common.ICancellationContext,
	connectionManager goConnectionManager.IService,
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
		CancellationContext,
		connectionManager,
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
