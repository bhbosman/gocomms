package impl

import (
	"context"
	"github.com/bhbosman/gocomms/connectionManager"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/gocomms/stacks/Bottom"
	"github.com/bhbosman/gocomms/stacks/KillConnection"
	"github.com/bhbosman/gocomms/stacks/Top"
	"github.com/bhbosman/gocomms/stacks/bvisMessageBreaker"
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

func CreateEmptyStack(
	connectionType internal.ConnectionType,
	_ string,
	_ interface{},
	_ connectionManager.IConnectionManager,
	_ context.Context,
	_ internal.CancelFunc,
	opts ...rxgo.Option) (*internal.TwoWayPipeDefinition, error) {
	result := internal.NewTwoWayPipeDefinition(nil)

	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Top.StackDefinition(connectionType, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Bottom.StackDefinition(connectionType, opts...)
	})
	return result, nil
}

func CreateWebSocketStack(
	connectionType internal.ConnectionType,
	connectionId string,
	userContext interface{},
	connectionManager connectionManager.IConnectionManager,
	cancelContext context.Context,
	stackCancelFunc internal.CancelFunc,
	opts ...rxgo.Option) (*internal.TwoWayPipeDefinition, error) {
	result := internal.NewTwoWayPipeDefinition(nil)

	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Top.StackDefinition(connectionType, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return websocket.StackDefinition(connectionType, stackCancelFunc, connectionManager, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Bottom.StackDefinition(connectionType, opts...)
	})
	return result, nil
}

func CreateDefaultStack(
	connectionType internal.ConnectionType,
	_ string,
	userContext interface{},
	connectionManager connectionManager.IConnectionManager,
	cancelContext context.Context,
	stackCancelFunc internal.CancelFunc,
	opts ...rxgo.Option) (*internal.TwoWayPipeDefinition, error) {
	result := internal.NewTwoWayPipeDefinition(nil)

	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Top.StackDefinition(connectionType, opts...)
	})

	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return KillConnection.StackDefinition(connectionType, cancelContext, stackCancelFunc, connectionManager, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Bottom.StackDefinition(connectionType, opts...)
	})
	return result, nil
}

func CreateCompressedStack(
	connectionType internal.ConnectionType,
	connectionId string,
	userContext interface{},
	connectionManager connectionManager.IConnectionManager,
	_ context.Context,
	stackCancelFunc internal.CancelFunc,
	opts ...rxgo.Option) (*internal.TwoWayPipeDefinition, error) {
	result := internal.NewTwoWayPipeDefinition(nil)
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Top.StackDefinition(connectionType, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return pingPong.StackDefinition(connectionType, connectionId, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return messageCompressor.StackDefinition(stackCancelFunc, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return messageNumber.StackDefinition(connectionType, userContext, stackCancelFunc, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return bvisMessageBreaker.StackDefinition(connectionType, connectionId, stackCancelFunc, nil, connectionManager, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Bottom.StackDefinition(connectionType, opts...)
	})
	return result, nil
}

func CreateUnCompressedStack(
	connectionType internal.ConnectionType,
	connectionId string,
	userContext interface{},
	connectionManager connectionManager.IConnectionManager,
	_ context.Context,
	stackCancelFunc internal.CancelFunc,
	opts ...rxgo.Option) (*internal.TwoWayPipeDefinition, error) {
	result := internal.NewTwoWayPipeDefinition(nil)
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Top.StackDefinition(connectionType, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return pingPong.StackDefinition(connectionType, connectionId, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return messageNumber.StackDefinition(connectionType, userContext, stackCancelFunc, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return bvisMessageBreaker.StackDefinition(connectionType, connectionId, stackCancelFunc, nil, connectionManager, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Bottom.StackDefinition(connectionType, opts...)
	})
	return result, nil
}

func CreateCompressedTlsStack(
	connectionType internal.ConnectionType,
	connectionId string,
	userContext interface{},
	connectionManager connectionManager.IConnectionManager,
	_ context.Context,
	stackCancelFunc internal.CancelFunc,
	opts ...rxgo.Option) (*internal.TwoWayPipeDefinition, error) {
	result := internal.NewTwoWayPipeDefinition(nil)
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Top.StackDefinition(connectionType, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return pingPong.StackDefinition(connectionType, connectionId, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return messageCompressor.StackDefinition(stackCancelFunc, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return messageNumber.StackDefinition(connectionType, userContext, stackCancelFunc, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return bvisMessageBreaker.StackDefinition(connectionType, connectionId, stackCancelFunc, nil, connectionManager, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return tlsConnection.StackDefinition(connectionType, stackCancelFunc, connectionManager, connectionId, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Bottom.StackDefinition(connectionType, opts...)
	})
	return result, nil
}

func CreateUnCompressedTlsStack(
	connectionType internal.ConnectionType,
	connectionId string,
	userContext interface{},
	connectionManager connectionManager.IConnectionManager,
	_ context.Context,
	stackCancelFunc internal.CancelFunc,
	opts ...rxgo.Option) (*internal.TwoWayPipeDefinition, error) {
	result := internal.NewTwoWayPipeDefinition(nil)
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Top.StackDefinition(connectionType, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return pingPong.StackDefinition(connectionType, connectionId, opts...)
	})

	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return messageNumber.StackDefinition(connectionType, userContext, stackCancelFunc, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return bvisMessageBreaker.StackDefinition(connectionType, connectionId, stackCancelFunc, nil, connectionManager, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return tlsConnection.StackDefinition(connectionType, stackCancelFunc, connectionManager, connectionId, opts...)
	})
	result.AddStackDefinitionFunc(func() (*internal.StackDefinition, error) {
		return Bottom.StackDefinition(connectionType, opts...)
	})
	return result, nil
}

const TransportFactoryCompressedTlsName = "CompressedTLS"
const TransportFactoryUnCompressedTlsName = "UncompressedTLS"
const TransportFactoryCompressedName = "Compressed"
const TransportFactoryUnCompressedName = "Uncompressed"
const TransportFactoryEmptyName = "Empty"
const WebSocketName = "WebSocket"

func GetNamedStack(name string) (TransportFactoryFunction, error) {
	switch name {
	case TransportFactoryCompressedTlsName:
		return CreateCompressedTlsStack, nil
	case TransportFactoryUnCompressedTlsName:
		return CreateUnCompressedTlsStack, nil
	case TransportFactoryCompressedName:
		return CreateCompressedStack, nil
	case TransportFactoryUnCompressedName:
		return CreateUnCompressedStack, nil
	case TransportFactoryEmptyName:
		return CreateEmptyStack, nil
	case WebSocketName:
		return CreateWebSocketStack, nil
	default:
		return CreateDefaultStack, nil
	}
}
