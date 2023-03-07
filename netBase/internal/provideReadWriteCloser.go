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
			Target: func() (net.Conn, io.ReadWriteCloser, io.Reader, io.Writer, io.Closer, error) {
				if conn == nil {
					return nil, nil, nil, nil, nil, fmt.Errorf("connection is nil. Please resolve")
				}
				return conn, conn, conn, conn, conn, nil
			},
		},
	)
}
