package tlsConnection

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/gomessageblock"
	"github.com/bhbosman/goprotoextra"
	"sync"

	"github.com/bhbosman/gocomms/stacks/internal/connectionWrapper"
	"github.com/reactivex/rxgo/v2"

	"go.uber.org/multierr"
	"io"
	"net"
)

func StackDefinition(
	connectionType internal.ConnectionType,
	stackCancelFunc internal.CancelFunc,
	connectionManager rxgo.IPublishToConnectionManager,
	connectionId string,
	opts ...rxgo.Option) (*internal.StackDefinition, error) {
	if stackCancelFunc == nil {
		return nil, goerrors.InvalidParam
	}

	nextOutBoundPath := func(ctx context.Context, nextOutboundChannel *internal.ChannelManager) connectionWrapper.ConnWrapperNext {
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

	// globals
	var connWrapper *connectionWrapper.ConnWrapper
	var pipeWriteClose io.WriteCloser
	var upgradedConnection net.Conn
	var nextInBoundChannel, nextOutboundChannel *internal.ChannelManager
	const stackName = "TLS"
	var stackIndex int
	// wg is here to make sure, that the upgradedConnection is properly assigned, before data is read/write to it
	upgradedConnectionAssignedWaitGroup := sync.WaitGroup{}
	upgradedConnectionAssignedWaitGroup.Add(1)
	return &internal.StackDefinition{
		Name: stackName,
		Inbound: func(inOutBoundParams internal.InOutBoundParams) internal.BoundDefinition {
			nextInBoundChannel = internal.NewChannelManager(make(chan rxgo.Item), "inbound TlsConnection", connectionId)
			stackIndex = inOutBoundParams.Index
			return internal.BoundDefinition{
				PipeDefinition: func(params internal.PipeDefinitionParams) (rxgo.Observable, error) {
					if stackCancelFunc == nil {
						return nil, goerrors.InvalidParam
					}
					_ = params.Obs.(rxgo.InOutBoundObservable).DoOnNextInOutBound(
						inOutBoundParams.Index,
						params.ConnectionId,
						stackName,
						rxgo.StreamDirectionInbound,
						connectionManager,
						func(ctx context.Context, rws goprotoextra.ReadWriterSize) {
							upgradedConnectionAssignedWaitGroup.Wait()
							_, err := io.Copy(pipeWriteClose, rws)
							if err != nil {
								return
							}
						}, opts...)
					nextObs := rxgo.FromChannel(nextInBoundChannel.Items)
					return nextObs, nil
				},
				PipeState: internal.PipeState{
					Start: func(ctx context.Context) error {
						return ctx.Err()
					},
					End: func() error {
						return nextInBoundChannel.Close()
					},
				},
			}
		},
		Outbound: func(inOutBoundParams internal.InOutBoundParams) internal.BoundDefinition {
			nextOutboundChannel = internal.NewChannelManager(make(chan rxgo.Item), "outbound Tls Connection", connectionId)
			return internal.BoundDefinition{
				PipeDefinition: func(params internal.PipeDefinitionParams) (rxgo.Observable, error) {
					if stackCancelFunc == nil {
						return nil, goerrors.InvalidParam
					}
					_ = params.Obs.(rxgo.InOutBoundObservable).DoOnNextInOutBound(
						inOutBoundParams.Index,
						params.ConnectionId,
						stackName,
						rxgo.StreamDirectionOutbound,
						connectionManager,
						func(ctx context.Context, size goprotoextra.ReadWriterSize) {
							upgradedConnectionAssignedWaitGroup.Wait()
							_, err := io.Copy(upgradedConnection, size)
							if err != nil {
								stackCancelFunc("copy data to upgradedConnection", false, err)
							}
						}, opts...)
					nextObs := rxgo.FromChannel(nextOutboundChannel.Items, opts...)
					return nextObs, nil
				},
				PipeState: internal.PipeState{
					Start: func(ctx context.Context) error {
						return ctx.Err()
					},
					End: func() error {
						return nextOutboundChannel.Close()
					},
				},
			}
		},
		StackState: internal.StackState{
			Start: func(startParams internal.StackStartStateParams) (net.Conn, error) {
				if startParams.Ctx.Err() != nil {
					return nil, startParams.Ctx.Err()
				}
				var pipeRead io.Reader
				pipeRead, pipeWriteClose = internal.Pipe(startParams.Ctx)
				if startParams.Ctx.Err() != nil {
					return nil, startParams.Ctx.Err()
				}
				connWrapper = connectionWrapper.NewConnWrapper(
					startParams.Conn,
					startParams.Ctx,
					pipeRead,
					nextOutBoundPath(startParams.Ctx, nextOutboundChannel))

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
				upgradedConnection = tlsConn
				upgradedConnectionAssignedWaitGroup.Done()
				if startParams.Ctx.Err() != nil {
					return nil, startParams.Ctx.Err()
				}

				go internal.ReadDataFromConnection(
					upgradedConnection,
					stackCancelFunc,
					startParams.Ctx,
					connectionManager,
					connectionId,
					stackIndex-1,
					"Read TLS Connection",
					func(rws goprotoextra.IReadWriterSize, cancelCtx context.Context, CancelFunc internal.CancelFunc) {
						nextInBoundChannel.Send(cancelCtx, rws)
					})

				return startParams.Conn, startParams.Ctx.Err()
			},
			End: func(endParams internal.StackEndStateParams) error {
				err := pipeWriteClose.Close()
				err = multierr.Append(err, upgradedConnection.Close())
				err = multierr.Append(err, connWrapper.Close())
				return err
			},
		},
	}, nil
}
