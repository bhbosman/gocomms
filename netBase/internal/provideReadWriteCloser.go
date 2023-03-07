package internal

import (
	"fmt"
	"go.uber.org/fx"
	"io"
)

func ProvideReadWriteCloser(rwc io.ReadWriteCloser) fx.Option {
	return fx.Provide(
		fx.Annotated{
			//Name: "Primary",
			Target: func() (io.ReadWriteCloser, io.Reader, io.Writer, io.Closer, error) {
				if rwc == nil {
					return nil, nil, nil, nil, fmt.Errorf("connection is nil. Please resolve")
				}
				return rwc, rwc, rwc, rwc, nil
			},
		},
	)
}
