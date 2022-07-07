package common

import "context"

type OnCreateConnectionDefault struct {
}

func (self *OnCreateConnectionDefault) OnCreateConnection(reference string, err error, ctx context.Context, cancelFunc context.CancelFunc) {
	//no nothing
	return
}
