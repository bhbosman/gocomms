package impl

import (
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/goerrors"
	"go.uber.org/fx"
)

type ConnectionReactorFactories struct {
	m map[string]intf.IConnectionReactorFactory
}

//func (self *ConnectionReactorFactories) CreateClientContext(
//	logger *gologging.SubSystemLogger,
//	ConnectionName string,
//	cancelCtx context.Context,
//	cancelFunc context.CancelFunc,
//	name string,
//	userContext interface{}) (intf.IConnectionReactor, error) {
//	if factory, ok := self.m[name]; ok {
//		return factory.Create(ConnectionName, cancelCtx, cancelFunc, logger, userContext), nil
//	}
//	return nil, goerrors.InvalidParam
//}

const ConnectionReactorFactoryConst = "ConnectionReactorFactory"

func newConnectionReactorFactories(
	params struct {
		fx.In
		Factories []intf.IConnectionReactorFactory `group:"ConnectionReactorFactory"`
	}) (*ConnectionReactorFactories, error) {
	m := make(map[string]intf.IConnectionReactorFactory)
	for _, ccf := range params.Factories {
		if _, ok := m[ccf.Name()]; ok {
			return nil, goerrors.DuplicateKey
		}
		m[ccf.Name()] = ccf
	}
	result := &ConnectionReactorFactories{
		m: m,
	}
	return result, nil
}
