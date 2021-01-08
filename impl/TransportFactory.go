package impl

import (
	"context"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/gocomms/stacks/Bottom"
	"github.com/bhbosman/gocomms/stacks/KillConnection"
	"github.com/bhbosman/gocomms/stacks/Top"
	"github.com/bhbosman/gocomms/stacks/messageBreaker"
	"github.com/bhbosman/gocomms/stacks/messageCompressor"
	"github.com/bhbosman/gocomms/stacks/messageNumber"
	"github.com/bhbosman/gocomms/stacks/pingPong"
	"github.com/bhbosman/gocomms/stacks/tlsConnection"
	"github.com/bhbosman/gocomms/stacks/websocket"
	"github.com/reactivex/rxgo/v2"
)

type TransportFactoryFunction func(
	connectionType internal.ConnectionType,
	connectionId string,
	userContext interface{},
	connectionManager connectionManager.IConnectionManager,
	cancelContext context.Context,
	cancelFunc internal.CancelFunc,
	opts ...rxgo.Option) (*internal.TwoWayPipeDefinition, error)

type TransportFactory struct {
}

func (self *TransportFactory) CreateEmpty(
	connectionType internal.ConnectionType,
	_ string,
	_ interface{},
	_ connectionManager.IConnectionManager,
	_ context.Context,
	_ internal.CancelFunc,
	opts ...rxgo.Option) (*internal.TwoWayPipeDefinition, error) {
	result := internal.NewTwoWayPipeDefinition(nil)

	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Top.StackDefinition()
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Bottom.StackDefinition()
	})
	return result, nil
}

func (self *TransportFactory) CreateWebSocket(
	connectionType internal.ConnectionType,
	connectionId string,
	userContext interface{},
	connectionManager connectionManager.IConnectionManager,
	cancelContext context.Context,
	stackCancelFunc internal.CancelFunc,
	opts ...rxgo.Option) (*internal.TwoWayPipeDefinition, error) {
	result := internal.NewTwoWayPipeDefinition(nil)

	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Top.StackDefinition()
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return websocket.StackDefinition(connectionType, stackCancelFunc, connectionManager, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Bottom.StackDefinition()
	})
	return result, nil
}

func (self *TransportFactory) CreateDefault(
	connectionType internal.ConnectionType,
	_ string,
	userContext interface{},
	connectionManager connectionManager.IConnectionManager,
	cancelContext context.Context,
	stackCancelFunc internal.CancelFunc,
	opts ...rxgo.Option) (*internal.TwoWayPipeDefinition, error) {
	result := internal.NewTwoWayPipeDefinition(nil)

	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Top.StackDefinition()
	})

	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return KillConnection.StackDefinition(cancelContext, stackCancelFunc, connectionManager, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Bottom.StackDefinition()
	})
	return result, nil
}

func (self *TransportFactory) CreateCompressed(
	connectionType internal.ConnectionType,
	connectionId string,
	userContext interface{},
	connectionManager connectionManager.IConnectionManager,
	_ context.Context,
	stackCancelFunc internal.CancelFunc,
	opts ...rxgo.Option) (*internal.TwoWayPipeDefinition, error) {
	result := internal.NewTwoWayPipeDefinition(nil)
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Top.StackDefinition()
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return pingPong.StackDefinition(connectionId, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return messageCompressor.StackDefinition(stackCancelFunc, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return messageNumber.StackDefinition(userContext, stackCancelFunc, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return messageBreaker.StackDefinition(connectionId, stackCancelFunc, nil, connectionManager, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Bottom.StackDefinition()
	})
	return result, nil
}

func (self *TransportFactory) CreateUnCompressed(
	connectionType internal.ConnectionType,
	connectionId string,
	userContext interface{},
	connectionManager connectionManager.IConnectionManager,
	_ context.Context,
	stackCancelFunc internal.CancelFunc,
	opts ...rxgo.Option) (*internal.TwoWayPipeDefinition, error) {
	result := internal.NewTwoWayPipeDefinition(nil)
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Top.StackDefinition()
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return pingPong.StackDefinition(connectionId, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return messageNumber.StackDefinition(userContext, stackCancelFunc, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return messageBreaker.StackDefinition(connectionId, stackCancelFunc, nil, connectionManager, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Bottom.StackDefinition()
	})
	return result, nil
}

func (self *TransportFactory) CreateCompressedTls(
	connectionType internal.ConnectionType,
	connectionId string,
	userContext interface{},
	connectionManager connectionManager.IConnectionManager,
	_ context.Context,
	stackCancelFunc internal.CancelFunc,
	opts ...rxgo.Option) (*internal.TwoWayPipeDefinition, error) {
	result := internal.NewTwoWayPipeDefinition(nil)
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Top.StackDefinition()
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return pingPong.StackDefinition(connectionId, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return messageCompressor.StackDefinition(stackCancelFunc, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return messageNumber.StackDefinition(userContext, stackCancelFunc, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return messageBreaker.StackDefinition(connectionId, stackCancelFunc, nil, connectionManager, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return tlsConnection.StackDefinition(connectionType, stackCancelFunc, connectionManager, connectionId)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Bottom.StackDefinition()
	})
	return result, nil
}

func (self *TransportFactory) CreateUnCompressedTls(
	connectionType internal.ConnectionType,
	connectionId string,
	userContext interface{},
	connectionManager connectionManager.IConnectionManager,
	_ context.Context,
	stackCancelFunc internal.CancelFunc,
	opts ...rxgo.Option) (*internal.TwoWayPipeDefinition, error) {
	result := internal.NewTwoWayPipeDefinition(nil)
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Top.StackDefinition()
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return pingPong.StackDefinition(connectionId, opts...)
	})

	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return messageNumber.StackDefinition(userContext, stackCancelFunc, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return messageBreaker.StackDefinition(connectionId, stackCancelFunc, nil, connectionManager, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return tlsConnection.StackDefinition(connectionType, stackCancelFunc, connectionManager, connectionId)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Bottom.StackDefinition()
	})
	return result, nil
}

const TransportFactoryCompressedTlsName = "CompressedTLS"
const TransportFactoryUnCompressedTlsName = "UncompressedTLS"
const TransportFactoryCompressedName = "Compressed"
const TransportFactoryUnCompressedName = "Uncompressed"
const TransportFactoryEmptyName = "Empty"
const WebSocketName = "WebSocket"

func (self *TransportFactory) Get(name string) (TransportFactoryFunction, error) {
	switch name {
	case TransportFactoryCompressedTlsName:
		return self.CreateCompressedTls, nil
	case TransportFactoryUnCompressedTlsName:
		return self.CreateUnCompressedTls, nil
	case TransportFactoryCompressedName:
		return self.CreateCompressed, nil
	case TransportFactoryUnCompressedName:
		return self.CreateUnCompressed, nil
	case TransportFactoryEmptyName:
		return self.CreateEmpty, nil
	case WebSocketName:
		return self.CreateWebSocket, nil
	default:
		return self.CreateDefault, nil
	}
}

func NewTransportFactory() *TransportFactory {
	result := &TransportFactory{}
	return result
}
