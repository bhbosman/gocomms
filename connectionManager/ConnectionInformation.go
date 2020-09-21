package connectionManager

import (
	"context"

	"github.com/icza/gox/fmtx"
	"github.com/reactivex/rxgo/v2"
	"sync"
	"time"
)

type StackPropertyKey struct {
	Index     int
	Name      string
	Direction rxgo.StreamDirection
}

type StackPropertyValue struct {
	msgCount  int
	byteCount int
}

func (self StackPropertyValue) MsgCount() string {
	return fmtx.FormatInt(int64(self.msgCount), 3, ',')
}

func (self StackPropertyValue) ByteCount() string {
	return fmtx.FormatSize(int64(self.byteCount), fmtx.SizeUnitAuto, 2)
}

type ConnectionInformation struct {
	Id              string
	CancelFunc      context.CancelFunc
	Name            string
	Status          string
	ConnectionTime  time.Time
	StackProperties map[StackPropertyKey]StackPropertyValue
	mutex           sync.Mutex
}

func NewConnectionInformation(id string, function context.CancelFunc) *ConnectionInformation {

	return &ConnectionInformation{
		Id:              id,
		CancelFunc:      function,
		StackProperties: make(map[StackPropertyKey]StackPropertyValue),
		ConnectionTime:  time.Now().Truncate(time.Second),
	}
}
