package common

type IStackDefinition interface {
	Name() string
	OnInbound() BoundResult
	OnOutbound() BoundResult
	StackState() IStackState
}
