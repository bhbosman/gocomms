package internal

import (
	"fmt"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/gologging"
	"go.uber.org/fx"
	"log"
)

func CreateTransportLayer00(
	params struct {
		fx.In
		Logger               *gologging.SubSystemLogger
		TwoWayPipeDefinition *internal.TwoWayPipeDefinition
		Lifecycle            fx.Lifecycle
	}) ([]*internal.StackState, error) {
	params.Logger.LogWithLevel(0, func(logger *log.Logger) {
		logger.Printf(fmt.Sprintf("createTransportLayer00..."))
	})
	return params.TwoWayPipeDefinition.BuildStackState()
}
