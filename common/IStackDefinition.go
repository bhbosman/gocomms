package common

type IStackDefinition interface {
	Name() string
	Inbound() IBoundResult
	Outbound() IBoundResult
	StackState() *StackState
}
