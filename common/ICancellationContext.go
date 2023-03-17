package common

import "github.com/bhbosman/goConn"

type ICancellationContext = goConn.ICancellationContext

//
//type ICancellationContext interface {
//	Add(connectionId string, f func()) (bool, error)
//	Remove(connectionId string) error
//	Cancel()
//	CancelWithError(err error)
//}
