package internal

import (
	"fmt"
	"go.uber.org/fx"
	"io"
)

func ProvideReadWriteCloser(rwc io.ReadWriteCloser) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Target: func() (io.ReadWriteCloser, error) {
				if rwc == nil {
					return nil, fmt.Errorf("connection is nil. Please resolve")
				}
				return rwc, nil
			},
		},
	)
}
