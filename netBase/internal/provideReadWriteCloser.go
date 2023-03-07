package internal

import (
	"fmt"
	"go.uber.org/fx"
	"io"
	"net"
)

func ProvideReadWriteCloser(conn net.Conn) fx.Option {
	return fx.Provide(
		fx.Annotated{
			Name: "PrimaryConnection",
			Target: func() (net.Conn, io.Closer, error) {
				if conn == nil {
					return nil, nil, fmt.Errorf("connection is nil. Please resolve")
				}
				return conn, conn, nil
			},
		},
	)
}
