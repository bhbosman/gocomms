package common

type IStackBoundFactory interface {
	GetPipeDefinition() PipeDefinition
	GetPipeState() *PipeState
}
