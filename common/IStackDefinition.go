package common

type IStackDefinition interface {
	GetName() string
	GetInbound() IBoundResult
	GetOutbound() IBoundResult
	GetStackState() *StackState
}
