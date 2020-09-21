package connectionManager

import (
	"context"
	"github.com/cskr/pubsub"
	"github.com/reactivex/rxgo/v2"
)

type IConnectionManagerWithData interface {
	IConnectionManager
	Clear()
	Publish()
}
type data struct {
	sub *pubsub.PubSub
	m   map[string]*ConnectionInformation
}

func (self *data) PublishStackData(index int, connectionId, name string, direction rxgo.StreamDirection, msgValue, byteValue int) {
	if ci, ok := self.m[connectionId]; ok {
		ci.mutex.Lock()
		defer ci.mutex.Unlock()
		data := StackPropertyValue{
			msgCount:  msgValue,
			byteCount: byteValue,
		}
		ci.StackProperties[StackPropertyKey{
			Index:     index,
			Name:      name,
			Direction: direction,
		}] = data
	}
}

func (self *data) CloseAllConnections(ctx context.Context) error {
	for _, connectionInformation := range self.m {
		go func(ci *ConnectionInformation) {
			ci.CancelFunc()
		}(connectionInformation)
	}
	return nil
}

func (self *data) CloseConnection(id string) {
	if ci, ok := self.m[id]; ok {
		ci.CancelFunc()
	}
}

func (self *data) GetConnections(ctx context.Context) ([]*ConnectionInformation, error) {
	var result []*ConnectionInformation
	for _, ci := range self.m {
		result = append(result, ci)
	}
	return result, nil
}

func (self *data) NameConnection(id string, name string) error {
	if ci, ok := self.m[id]; ok {
		ci.Name = name
	}
	return nil
}

func (self *data) StatusConnection(id string, status string) error {
	if ci, ok := self.m[id]; ok {
		ci.Status = status
	}
	return nil
}

func (self *data) Publish() {

}

func (self *data) Clear() {

}

func (self *data) RegisterConnection(id string, function context.CancelFunc) error {
	self.m[id] = NewConnectionInformation(id, function)
	return nil
}

func (self *data) DeregisterConnection(id string) error {
	delete(self.m, id)
	return nil
}

func (self *data) Start(ctx context.Context) error {
	return nil
}

func (self *data) Stop(ctx context.Context) error {
	return nil
}

func newData(sub *pubsub.PubSub) IConnectionManagerWithData {
	return &data{
		sub: sub,
		m:   make(map[string]*ConnectionInformation),
	}
}
