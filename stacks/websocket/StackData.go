package websocket

import (
	"context"
	"fmt"
	"github.com/bhbosman/gocommon/stream"
	"github.com/bhbosman/gocomms/internal"
	"github.com/bhbosman/gocomms/intf"
	"github.com/bhbosman/gocomms/stacks/internal/connectionWrapper"
	"github.com/bhbosman/gocomms/stacks/websocket/wsmsg"
	"github.com/bhbosman/goerrors"
	"github.com/bhbosman/gomessageblock"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"go.uber.org/multierr"
	"io"
	"net"
	"net/url"
	"reflect"
	"time"
)

func WrongStackDataError(connectionType internal.ConnectionType, stackData interface{}) error {
	return internal.NewWrongStackDataType(
		StackName,
		connectionType,
		reflect.TypeOf((*StackData)(nil)),
		reflect.TypeOf(stackData))
}

type StackData struct {
	upgradedConnection         net.Conn
	stackCancelFunc            internal.CancelFunc
	LastPongReceived           time.Time
	nextInBoundChannelManager  *internal.ChannelManager
	nextOutboundChannelManager *internal.ChannelManager
	tempStep                   *internal.ChannelManager

	connWrapper    *connectionWrapper.ConnWrapper
	pipeWriteClose io.WriteCloser
}

func (self *StackData) Close() error {
	var err error = nil
	err = multierr.Append(err, self.nextInBoundChannelManager.Close())
	err = multierr.Append(err, self.nextOutboundChannelManager.Close())
	err = multierr.Append(err, self.tempStep.Close())
	if self.upgradedConnection != nil {
		err = multierr.Append(err, self.upgradedConnection.Close())
	}
	if self.pipeWriteClose != nil {
		err = multierr.Append(err, self.pipeWriteClose.Close())
	}
	if self.connWrapper != nil {
		err = multierr.Append(err, self.connWrapper.Close())
	}
	return err
}

func nextOutBoundPath(
	ctx context.Context,
	nextOutboundChannelManager *internal.ChannelManager) connectionWrapper.ConnWrapperNext {
	return func(b []byte) (n int, err error) {
		if ctx.Err() != nil {
			return 0, ctx.Err()
		}
		rws := gomessageblock.NewReaderWriterSize(len(b))
		if ctx.Err() != nil {
			return 0, ctx.Err()
		}
		n, err = rws.Write(b)
		if err != nil {
			return 0, err
		}
		if ctx.Err() != nil {
			return 0, ctx.Err()
		}
		nextOutboundChannelManager.Send(ctx, rws)
		return n, nil
	}
}

func (self *StackData) sendPing(ctx context.Context) {
	msg := &wsmsg.WebSocketMessage{
		OpCode: wsmsg.WebSocketMessage_OpPing,
	}
	marshall, err := stream.Marshall(msg)
	if err != nil {
		return
	}
	self.tempStep.Send(ctx, marshall)
}

func (self *StackData) triggerPingLoop(cancelContext context.Context) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for true {
		select {
		case <-cancelContext.Done():
			return
		case <-ticker.C:
			if time.Now().Sub(self.LastPongReceived) > time.Second*10 {
				self.stackCancelFunc("pong time out", true, goerrors.TimeOut)
				return
			}
			self.sendPing(cancelContext)
		}
	}
}

func (self *StackData) OnStart(startParams internal.StackStartStateParams) (net.Conn, error) {
	if startParams.Ctx.Err() != nil {
		return nil, startParams.Ctx.Err()
	}

	// create map and fill it in with some values
	inputValues := make(map[string]interface{})
	inputValues["url"] = startParams.Url
	inputValues["localAddr"] = startParams.Conn.LocalAddr()
	inputValues["remoteAddr"] = startParams.Conn.RemoteAddr()
	outputValues, err := startParams.ConnectionReactorFactory.Values(inputValues)
	if err != nil {
		return nil, err
	}

	// get header information from outputValues
	header := make(ws.HandshakeHeaderHTTP)
	if additionalHeaderInformation, ok := outputValues["connectionHeader"]; ok {
		if connectionHeader, isMap := additionalHeaderInformation.(map[string][]string); isMap {
			for k, v := range connectionHeader {
				header[k] = v
			}
		}
	}

	// build websocket dialer that will be used
	dialer := ws.Dialer{
		Header: header,
		NetDial: func(ctx context.Context, network, addr string) (net.Conn, error) {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			return self.connWrapper, nil
		},
	}
	// On error exit
	if startParams.Ctx.Err() != nil {
		return nil, startParams.Ctx.Err()
	}
	self.upgradedConnection, _, _, err = dialer.Dial(startParams.Ctx, startParams.Url.String())

	// On error exit
	if err != nil {
		return nil, err
	}

	// On error exit
	if startParams.Ctx.Err() != nil {
		return nil, startParams.Ctx.Err()
	}

	go self.connectionLoop(startParams.Ctx)
	go self.triggerPingLoop(startParams.Ctx)
	return self.upgradedConnection, startParams.Ctx.Err()
}

func (self *StackData) connectionLoop(ctx context.Context) {
	sendMessage := func(message *wsmsg.WebSocketMessage) {
		stm, err := stream.Marshall(message)
		if err != nil {
			//return
		}
		self.nextInBoundChannelManager.Send(ctx, stm)
	}
	message := wsmsg.WebSocketMessage{
		OpCode:  wsmsg.WebSocketMessage_OpStartLoop,
		Message: nil,
	}
	sendMessage(&message)

	for {
		msgs, err := wsutil.ReadServerMessage(self.upgradedConnection, nil)
		if err != nil {
			return
		}
		if ctx.Err() != nil {
			return
		}
		for _, msg := range msgs {
			if ctx.Err() != nil {
				return
			}
			var message wsmsg.WebSocketMessage
			switch msg.OpCode {
			case ws.OpContinuation:
				message = wsmsg.WebSocketMessage{
					OpCode:  wsmsg.WebSocketMessage_OpContinuation,
					Message: msg.Payload,
				}
			case ws.OpText:
				message = wsmsg.WebSocketMessage{
					OpCode:  wsmsg.WebSocketMessage_OpText,
					Message: msg.Payload,
				}
			case ws.OpBinary:
				message = wsmsg.WebSocketMessage{
					OpCode:  wsmsg.WebSocketMessage_OpBinary,
					Message: msg.Payload,
				}
			case ws.OpClose:
				self.stackCancelFunc(
					"close message received",
					true,
					fmt.Errorf(string(msg.Payload)))
			case ws.OpPing:
				print("+")
				err := wsutil.WriteClientMessage(self.upgradedConnection, ws.OpPong, msg.Payload)
				if err != nil {
					self.stackCancelFunc("creating pong payload", true, err)
					return
				}
				continue
			case ws.OpPong:
				_ = self.LastPongReceived.UnmarshalBinary(msg.Payload)
				continue
			default:
				continue
			}
			if ctx.Err() != nil {
				return
			}
			sendMessage(&message)
			message.Reset()
		}
	}
}

func NewStackData(conn net.Conn, stackCancelFunc internal.CancelFunc, ctx context.Context, Url *url.URL, cfr intf.IConnectionReactorFactoryExtractValues) (*StackData, error) {
	nextInBoundChannelManager := internal.NewChannelManager("", "")
	nextOutboundChannelManager := internal.NewChannelManager("", "")
	tempStep := internal.NewChannelManager("", "")

	var pipeRead io.Reader
	var pipeWriteCloser io.WriteCloser
	pipeRead, pipeWriteCloser = internal.Pipe(ctx)

	connWrapper := connectionWrapper.NewConnWrapper(
		conn,
		ctx,
		pipeRead,
		nextOutBoundPath(ctx, nextOutboundChannelManager))

	return &StackData{
		upgradedConnection:         nil,
		stackCancelFunc:            stackCancelFunc,
		LastPongReceived:           time.Now(),
		nextInBoundChannelManager:  nextInBoundChannelManager,
		nextOutboundChannelManager: nextOutboundChannelManager,
		tempStep:                   tempStep,
		connWrapper:                connWrapper,
		pipeWriteClose:             pipeWriteCloser,
	}, nil
}
