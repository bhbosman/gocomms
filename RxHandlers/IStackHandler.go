package RxHandlers

type IStackHandler interface {
	ReadMessage(interface{}) (interface{}, bool, error)
}
