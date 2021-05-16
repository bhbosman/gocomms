package internal

import (
	"context"
	"fmt"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gologging"
	"go.uber.org/fx"
	"log"
)

func CreateClientContext(
	params struct {
		fx.In
		Lifecycle      fx.Lifecycle
		CancelCtx      context.Context
		CancelFunc     context.CancelFunc
		ConnectionName string `name:"ConnectionName"`
		Logger         *gologging.SubSystemLogger
		Cfr            intf.IConnectionReactorFactory
		ClientContext  interface{} `name:"UserContext"`
	}) (intf.IConnectionReactor, error) {
	params.Logger.LogWithLevel(0, func(logger *log.Logger) {
		logger.Printf(fmt.Sprintf("createTransportLayer..."))
	})
	return params.Cfr.Create(
		params.Cfr.Name(),
		params.CancelCtx,
		params.CancelFunc,
		params.Logger,
		params.ClientContext), nil
}
