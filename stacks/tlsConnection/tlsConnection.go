package tlsConnection

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/gocomms/stacks/internal/connectionWrapper"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"github.com/google/uuid"
	"github.com/reactivex/rxgo/v2"
	"net/url"
	"reflect"

	"go.uber.org/multierr"
	"io"
	"net"
)

const StackName = "TLS"

type StackData struct {
	connectionType      internal.ConnectionType
	Conn                io.ReadWriter
	connWrapper         *connectionWrapper.ConnWrapper
	pipeWriteClose      io.WriteCloser
	nextOutboundChannel *internal.ChannelManager
	nextInBoundChannel  *internal.ChannelManager
	upgradedConnection  net.Conn
	stackIndex          int
}

func NewStackData(connectionType internal.ConnectionType, Conn net.Conn, connectionId string, ctx context.Context) (*StackData, error) {
	nextOutboundChannel := internal.NewChannelManager("outbound Tls Connection", connectionId)
	nextInBoundChannel := internal.NewChannelManager("inbound TlsConnection", connectionId)

	var tempPipeRead io.Reader
	var tempPipeWriteClose io.WriteCloser
	tempPipeRead, tempPipeWriteClose = internal.Pipe(ctx)
	pipeWriteClose := tempPipeWriteClose
	connWrapper := connectionWrapper.NewConnWrapper(
		Conn,
		ctx,
		tempPipeRead,
		nextOutBoundPath(ctx, nextOutboundChannel))

	var tlsConn *tls.Conn
	if connectionType == internal.ServerConnection {
		cer, err := tls.X509KeyPair(_serverPem, _serverKey)
		if err != nil {
			return nil, err
		}
		config := &tls.Config{
			ServerName:   "127.0.0.1",
			Certificates: []tls.Certificate{cer},
			VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
				return nil
			},
			VerifyConnection: func(state tls.ConnectionState) error {
				return nil
			},
		}
		tlsConn = tls.Server(connWrapper, config)
	} else {
		config := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         "localhost"}
		tlsConn = tls.Client(connWrapper, config)
	}
	upgradedConnection := tlsConn

	return &StackData{
		connectionType:      connectionType,
		Conn:                Conn,
		connWrapper:         connWrapper,
		pipeWriteClose:      pipeWriteClose,
		nextOutboundChannel: nextOutboundChannel,
		nextInBoundChannel:  nextInBoundChannel,
		upgradedConnection:  upgradedConnection,
	}, nil
}

func (self *StackData) Close() error {
	var err error = nil
	err = multierr.Append(err, self.connWrapper.Close())
	err = multierr.Append(err, self.pipeWriteClose.Close())
	err = multierr.Append(err, self.nextOutboundChannel.Close())
	err = multierr.Append(err, self.nextInBoundChannel.Close())
	err = multierr.Append(err, self.upgradedConnection.Close())
	return nil
}

func nextOutBoundPath(ctx context.Context, nextOutboundChannel *internal.ChannelManager) connectionWrapper.ConnWrapperNext {
	return func(b []byte) (n int, err error) {
		if ctx.Err() != nil {
			return 0, ctx.Err()
		}
		dataToConnection := gomessageblock.NewReaderWriterSize(len(b))
		if ctx.Err() != nil {
			return 0, ctx.Err()
		}
		n, err = dataToConnection.Write(b)
		if err != nil {
			return 0, err
		}
		if ctx.Err() != nil {
			return 0, ctx.Err()
		}
		nextOutboundChannel.Send(ctx, dataToConnection)
		return n, nil
	}
}

func WrongStackDataError(connectionType internal.ConnectionType, stackData interface{}) error {
	return internal.NewWrongStackDataType(
		StackName,
		connectionType,
		reflect.TypeOf((*StackData)(nil)),
		reflect.TypeOf(stackData))
}

func StackDefinition(
	connectionType internal.ConnectionType,
	stackCancelFunc internal.CancelFunc,
	connectionManager rxgo.IPublishToConnectionManager,
	connectionId string,
	opts ...rxgo.Option) (*internal.StackDefinition, error) {

	if stackCancelFunc == nil {
		return nil, goerrors.InvalidParam
	}
	id := uuid.New()
	return &internal.StackDefinition{
		IId:  id,
		Name: StackName,
		Inbound: internal.NewBoundResultImpl(func(inOutBoundParams internal.InOutBoundParams) (internal.IStackBoundDefinition, error) {
			return &internal.StackBoundDefinition{
				PipeDefinition: func(stackData, pipeData interface{}, pipeParams internal.PipeDefinitionParams) (uuid.UUID, rxgo.Observable, error) {
					if stackCancelFunc == nil {
						return uuid.Nil, nil, goerrors.InvalidParam
					}
					StackData, ok := stackData.(*StackData)
					if !ok {
						return uuid.Nil, nil, WrongStackDataError(connectionType, stackData)
					}
					StackData.stackIndex = inOutBoundParams.Index

					_ = pipeParams.Obs.(rxgo.InOutBoundObservable).DoOnNextInOutBound(
						inOutBoundParams.Index,
						pipeParams.ConnectionId,
						StackName,
						rxgo.StreamDirectionInbound,
						connectionManager,
						func(ctx context.Context, rws goprotoextra.ReadWriterSize) {
							//upgradedConnectionAssignedWaitGroup.Wait()
							_, err := io.Copy(StackData.pipeWriteClose, rws)
							if err != nil {
								return
							}
						}, opts...)
					nextObs := rxgo.FromChannel(StackData.nextInBoundChannel.Items, opts...)
					return id, nextObs, nil
				},
			}, nil
		}),
		Outbound: internal.NewBoundResultImpl(func(inOutBoundParams internal.InOutBoundParams) (internal.IStackBoundDefinition, error) {
			return &internal.StackBoundDefinition{
				PipeDefinition: func(stackData, pipeData interface{}, pipeParams internal.PipeDefinitionParams) (uuid.UUID, rxgo.Observable, error) {
					StackData, ok := stackData.(*StackData)
					if !ok {
						return uuid.Nil, nil, WrongStackDataError(connectionType, stackData)
					}
					if stackCancelFunc == nil {
						return uuid.Nil, nil, goerrors.InvalidParam
					}
					_ = pipeParams.Obs.(rxgo.InOutBoundObservable).DoOnNextInOutBound(
						inOutBoundParams.Index,
						pipeParams.ConnectionId,
						StackName,
						rxgo.StreamDirectionOutbound,
						connectionManager,
						func(ctx context.Context, size goprotoextra.ReadWriterSize) {
							//upgradedConnectionAssignedWaitGroup.Wait()
							_, err := io.Copy(StackData.upgradedConnection, size)
							if err != nil {
								stackCancelFunc("copy data to upgradedConnection", false, err)
							}
						}, opts...)
					nextObs := rxgo.FromChannel(StackData.nextOutboundChannel.Items, opts...)
					return id, nextObs, nil
				},
			}, nil
		}),
		StackState: &internal.StackState{
			Id: id,
			Create: func(Conn net.Conn, Url *url.URL, ctx context.Context, CancelFunc internal.CancelFunc) (interface{}, error) {
				return NewStackData(connectionType, Conn, connectionId, ctx)
			},
			Destroy: func(stackData interface{}) error {
				if _, ok := stackData.(*StackData); !ok {
					return WrongStackDataError(connectionType, stackData)
				}
				if closer, ok := stackData.(io.Closer); ok {
					return closer.Close()
				}
				return nil
			},
			Start: func(stackData interface{}, startParams internal.StackStartStateParams) (net.Conn, error) {
				StackData, ok := stackData.(*StackData)
				if !ok {
					return nil, WrongStackDataError(connectionType, stackData)
				}
				go internal.ReadDataFromConnection(
					StackData.upgradedConnection,
					stackCancelFunc,
					startParams.Ctx,
					connectionManager,
					connectionId,
					StackData.stackIndex-1,
					"Read TLS Connection",
					func(rws goprotoextra.IReadWriterSize, cancelCtx context.Context, CancelFunc internal.CancelFunc) {
						StackData.nextInBoundChannel.Send(cancelCtx, rws)
					})

				return StackData.upgradedConnection, startParams.Ctx.Err()
			},
			Stop: func(stackData interface{}, endParams internal.StackEndStateParams) error {
				var err error = nil
				return err
			},
		},
	}, nil
}
