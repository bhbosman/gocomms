package common

type stackDefinition struct {
	name       string
	inbound    BoundResult
	outbound   BoundResult
	stackState *StackState
}

func NewStackDefinition(
	name string,
	inbound BoundResult,
	outbound BoundResult,
	stackState *StackState,
) (IStackDefinition, error) {
	return &stackDefinition{
		name:       name,
		inbound:    inbound,
		outbound:   outbound,
		stackState: stackState,
	}, nil
}

func (self *stackDefinition) Name() string {
	return self.name
}

func (self *stackDefinition) Inbound() BoundResult {
	return self.inbound
}

func (self *stackDefinition) Outbound() BoundResult {
	return self.outbound
}

func (self *stackDefinition) StackState() *StackState {
	return self.stackState
}
