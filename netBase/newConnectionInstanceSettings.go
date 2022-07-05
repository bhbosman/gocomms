package netBase

import "go.uber.org/fx"

type newConnectionInstanceSettings struct {
	options                         []fx.Option
	ProvideCreateIConnectionReactor fx.Option // func () (intf.IConnectionReactor, error){}
}
