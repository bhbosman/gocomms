package internal

import (
	"context"
	"fmt"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/gologging"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"log"
)

func CreateTransportLayer02Step01(
	params struct {
		fx.In
		ConnectionId         string `name:"ConnectionId"`
		ConnectionManager    connectionManager.IConnectionManager
		Logger               *gologging.SubSystemLogger
		CancelCtx            context.Context
		StackCancelFunc      internal.CancelFunc
		TwoWayPipeDefinition *internal.TwoWayPipeDefinition
		StackData            map[uuid.UUID]interface{} `name:"StackData"`
	}) (*internal.IncomingObs, error) {
	params.Logger.LogWithLevel(0, func(logger *log.Logger) {
		logger.Printf(fmt.Sprintf("createTransportLayer..."))
	})
	return params.TwoWayPipeDefinition.BuildIncomingObs(
		params.StackData,
		params.ConnectionId,
		params.ConnectionManager,
		params.CancelCtx,
		params.StackCancelFunc)
}
