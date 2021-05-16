package internal

import "net"

func ProvideNetConnection(conn net.Conn) func() (net.Conn, error) {
	return func() (net.Conn, error) {
		return conn, nil
	}
}
