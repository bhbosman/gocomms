package common

type stackBoundDefinition struct {
	PipeDefinition PipeDefinition
	PipeState      *PipeState
}

func (self *stackBoundDefinition) GetPipeDefinition() PipeDefinition {
	return self.PipeDefinition
}

func (self *stackBoundDefinition) GetPipeState() *PipeState {
	return self.PipeState
}

func NewBoundDefinition(pipeDefinition PipeDefinition, pipeState *PipeState) IStackBoundDefinition {
	return &stackBoundDefinition{
		PipeDefinition: pipeDefinition,
		PipeState:      pipeState,
	}
}
