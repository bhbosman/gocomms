package internal

type IStackBoundDefinition interface {
	GetPipeDefinition() PipeDefinition
	GetPipeState() *PipeState
}

type StackBoundDefinition struct {
	PipeDefinition PipeDefinition
	PipeState      *PipeState
}

func (self *StackBoundDefinition) GetPipeDefinition() PipeDefinition {
	return self.PipeDefinition
}

func (self *StackBoundDefinition) GetPipeState() *PipeState {
	return self.PipeState
}

func NewBoundDefinition(pipeDefinition PipeDefinition, pipeState *PipeState) IStackBoundDefinition {
	return &StackBoundDefinition{
		PipeDefinition: pipeDefinition,
		PipeState:      pipeState,
	}
}
