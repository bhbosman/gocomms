package common

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"io"
	"sync"
)

type cancellationContext struct {
	mutex         sync.Mutex
	cancelFunc    context.CancelFunc
	cancelContext context.Context
	logger        *zap.Logger
	f             map[string]func()
	cancelCalled  bool
	closer        io.Closer
	name          string
}

func (self *cancellationContext) Remove(connectionId string) error {
	if !self.cancelCalled {
		self.mutex.Lock()
		defer self.mutex.Unlock()
		delete(self.f, connectionId)
	}
	return nil
}

func (self *cancellationContext) CancelWithError(s string, err error) {
	self.Cancel(s)
}

func (self *cancellationContext) Add(connectionId string, f func()) (bool, error) {
	if !self.cancelCalled {
		self.mutex.Lock()
		defer self.mutex.Unlock()
		//
		if foundFunction, ok := self.f[connectionId]; ok {
			foundFunction()
		}
		self.f[connectionId] = f
		return true, nil
	}
	f()
	return false, nil
}

func (self *cancellationContext) CancelContext() context.Context {
	return self.cancelContext
}

func (self *cancellationContext) CancelFunc() context.CancelFunc {
	return func() {
		// todo: fix
		self.Cancel("")
	}

}

func (self *cancellationContext) Cancel(s string) {
	self.mutex.Lock()
	b := self.cancelCalled
	self.cancelCalled = true
	self.mutex.Unlock()
	if !b {
		self.logger.Info(fmt.Sprintf("Cancel func for connection called"))
		self.cancelFunc()
		if self.closer != nil {
			self.closer.Close()
		}
		self.mutex.Lock()
		fArray := make([]func(), 0, len(self.f))
		for _, f := range self.f {
			fArray = append(fArray, f)
		}
		self.f = make(map[string]func())
		self.mutex.Unlock()
		for _, f := range fArray {
			f()
		}
	}
}

func NewCancellationContext(
	name string,
	cancelFunc context.CancelFunc,
	cancelContext context.Context,
	logger *zap.Logger,
	closer io.Closer,
) ICancellationContext {
	return &cancellationContext{
		name:          name,
		cancelFunc:    cancelFunc,
		cancelContext: cancelContext,
		logger:        logger,
		cancelCalled:  false,
		closer:        closer,
		f:             make(map[string]func()),
	}
}
