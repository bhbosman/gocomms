package RxHandlers

type IRxNextHandler interface {
	OnSendData(i interface{})
	OnError(err error)
	OnComplete()
}
