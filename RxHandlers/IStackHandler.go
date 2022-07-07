package RxHandlers

type IStackHandler interface {
	ReadMessage(interface{}) error
}
