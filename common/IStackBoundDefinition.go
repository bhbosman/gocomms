package common

type IStackBoundDefinition interface {
	GetPipeDefinition() PipeDefinition
	GetPipeState() *PipeState
}
