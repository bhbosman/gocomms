package internal

import (
	"context"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/gologging"
)

func CreateStackCancelFunc(cancelFunc context.CancelFunc, logger *gologging.SubSystemLogger) internal.CancelFunc {
	return func(context string, inbound bool, err error) {
		_ = logger.ErrorWithDescription(context, err)
		cancelFunc()
	}
}
