package common

import (
	"context"
	"github.com/bhbosman/goConnectionManager"
	"github.com/bhbosman/gocommon/GoFunctionCounter"
	"github.com/bhbosman/gocommon/Services/interfaces"
	"github.com/bhbosman/gocommon/messages"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type NetAppFuncInParams struct {
	fx.In
	ParentContext       context.Context `name:"Application"`
	ConnectionManager   goConnectionManager.IService
	ZapLogger           *zap.Logger
	UniqueSessionNumber interfaces.IUniqueReferenceService
	GoFunctionCounter   GoFunctionCounter.IService
}

func NewNetAppFuncInParams(
	parentContext context.Context,
	connectionManager goConnectionManager.IService,
	zapLogger *zap.Logger,
	uniqueSessionNumber interfaces.IUniqueReferenceService,
	goFunctionCounter GoFunctionCounter.IService,
) NetAppFuncInParams {
	return NetAppFuncInParams{
		ParentContext:       parentContext,
		ConnectionManager:   connectionManager,
		ZapLogger:           zapLogger,
		UniqueSessionNumber: uniqueSessionNumber,
		GoFunctionCounter:   goFunctionCounter,
	}
}

type NetAppFuncInParamsCallback func(params NetAppFuncInParams) messages.CreateAppCallback
