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

func NewStackBoundDefinition(
	pipeDefinition PipeDefinition,
	pipeState *PipeState,
) IStackBoundFactory {
	return &stackBoundDefinition{
		PipeDefinition: pipeDefinition,
		PipeState:      pipeState,
	}
}
