package internal

import (
	"context"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/goprotoextra"
	"go.uber.org/fx"
)

func CreateToConnectionFunc(
	params struct {
		fx.In
		OutgoingObs *internal.OutgoingObs
		CancelCtx   context.Context
	}) goprotoextra.ToConnectionFunc {
	return func(rw goprotoextra.ReadWriterSize) error {
		err := params.CancelCtx.Err()
		if err != nil {
			return err
		}
		return params.OutgoingObs.SendOutgoingData(rw)
	}
}
