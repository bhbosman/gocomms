package connectionManager

import (
	"context"
	"github.com/bhbosman/goerrors"
	"github.com/cskr/pubsub"
	"github.com/reactivex/rxgo/v2"
	"time"
)

type IRegisterToConnectionManager interface {
	RegisterConnection(id string, function context.CancelFunc, CancelContext context.Context) error
	DeregisterConnection(id string) error
	NameConnection(id string, name string) error
	StatusConnection(id string, name string) error
}

type IObtainConnectionManagerInformation interface {
	GetConnections(ctx context.Context) ([]*ConnectionInformation, error)
}

type ICommandsToConnectionManager interface {
	CloseConnection(id string)
	CloseAllConnections(ctx context.Context) error
}

type IConnectionManager interface {
	rxgo.IPublishToConnectionManager
	IRegisterToConnectionManager
	IObtainConnectionManagerInformation
	ICommandsToConnectionManager
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type operation int

const (
	start operation = iota
	stop
	register
	deregister
	name
	cleanData
	status
	get
	closeConnection
	closeAllConnection
	publishInt
)

type command struct {
	cmd       operation
	id        string
	function  context.CancelFunc
	ctx       context.Context
	s         string
	index     int
	i, i2     int
	ch        chan interface{}
	direction rxgo.StreamDirection
}

type Impl struct {
	cmdChannel chan command
	cancelCtx  context.Context
	sub        *pubsub.PubSub
}

func (self *Impl) PublishStackData(index int, connectionId, name string, direction rxgo.StreamDirection, msgValue, byteValue int) {
	self.cmdChannel <- command{cmd: publishInt, index: index, id: connectionId, s: name, direction: direction, i: msgValue, i2: byteValue}
}

func (self *Impl) CloseAllConnections(ctx context.Context) error {
	ch := make(chan interface{})
	defer close(ch)
	self.cmdChannel <- command{cmd: closeAllConnection, ch: ch, ctx: ctx}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case ans := <-ch:
		switch v := ans.(type) {
		case error:
			return v
		}
	}
	return nil
}

func (self *Impl) CloseConnection(id string) {
	self.cmdChannel <- command{cmd: closeConnection, id: id, s: id}
}

func (self *Impl) GetConnections(ctx context.Context) ([]*ConnectionInformation, error) {
	ch := make(chan interface{})
	defer close(ch)
	self.cmdChannel <- command{cmd: get, ch: ch, ctx: ctx}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case ans := <-ch:
		switch v := ans.(type) {
		case error:
			return nil, v
		case []*ConnectionInformation:
			return v, nil
		}
	}
	return nil, goerrors.InvalidParam
}

func (self *Impl) NameConnection(id string, connectionName string) error {
	self.cmdChannel <- command{cmd: name, id: id, s: connectionName}
	return nil
}

func (self *Impl) StatusConnection(id string, connectionStatus string) error {
	self.cmdChannel <- command{cmd: status, id: id, s: connectionStatus}
	return nil
}

func (self *Impl) Start(ctx context.Context) error {
	self.cmdChannel <- command{cmd: start, ctx: ctx}
	return nil
}

func (self *Impl) Stop(ctx context.Context) error {
	self.cmdChannel <- command{cmd: stop, ctx: ctx}
	return nil
}

func (self *Impl) RegisterConnection(id string, function context.CancelFunc, CancelContext context.Context) error {
	self.cmdChannel <- command{cmd: register, id: id, function: function, ctx: CancelContext}
	return nil
}

func (self *Impl) DeregisterConnection(id string) error {
	self.cmdChannel <- command{cmd: deregister, id: id}
	return nil
}

func (self *Impl) start() {
	data, _ := newData(self.sub)
	timer := time.NewTicker(time.Second * 5)
	defer timer.Stop()
loop:
	for {
		select {
		case _, ok := <-timer.C:
			if ok {
				go func() {
					self.cmdChannel <- command{cmd: cleanData}
				}()
			}
		case <-self.cancelCtx.Done():
			break loop
		case cmd, ok := <-self.cmdChannel:
			if ok {
				switch cmd.cmd {
				case cleanData:
					data.Cleanup()
					continue loop
				case stop:
					_ = data.Stop(cmd.ctx)
					break loop
				case start:
					_ = data.Start(cmd.ctx)
					continue loop
				case register:
					_ = data.RegisterConnection(cmd.id, cmd.function, cmd.ctx)
					continue loop
				case deregister:
					_ = data.DeregisterConnection(cmd.id)
					continue loop
				case name:
					_ = data.NameConnection(cmd.id, cmd.s)
					continue loop
				case status:
					_ = data.StatusConnection(cmd.id, cmd.s)
					continue loop
				//case publish:
				//	data.Publish()
				//	continue loop
				case get:
					connections, err := data.GetConnections(cmd.ctx)
					if err != nil {
						cmd.ch <- err
						continue loop
					}
					cmd.ch <- connections
					continue loop
				case closeConnection:
					data.CloseConnection(cmd.s)
					continue loop
				case closeAllConnection:
					err := data.CloseAllConnections(cmd.ctx)
					cmd.ch <- err
					continue loop
				case publishInt:
					data.PublishStackData(cmd.index, cmd.id, cmd.s, cmd.direction, cmd.i, cmd.i2)
					continue loop
				}
			}
		}
	}
	data.Clear()
	// flushing
	for range self.cmdChannel {
	}
}

func NewConnectionManager(sub *pubsub.PubSub, cancelCtx context.Context) IConnectionManager {
	result := &Impl{
		cmdChannel: make(chan command),
		cancelCtx:  cancelCtx,
		sub:        sub,
	}
	go result.start()
	return result
}
