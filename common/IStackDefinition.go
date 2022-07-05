package common

type IStackDefinition interface {
	GetId() string
	GetName() string
	GetInbound() IBoundResult
	GetOutbound() IBoundResult
	GetStackState() *StackState
}
