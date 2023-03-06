package common

type IStackDefinition interface {
	Name() string
	Inbound() BoundResult
	Outbound() BoundResult
	StackState() IStackState
}
