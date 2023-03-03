package common

type stackDefinition struct {
	name       string
	inbound    IBoundResult
	outbound   IBoundResult
	stackState *StackState
}

func NewStackDefinition(
	name string,
	inbound IBoundResult,
	outbound IBoundResult,
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

func (self *stackDefinition) Inbound() IBoundResult {
	return self.inbound
}

func (self *stackDefinition) Outbound() IBoundResult {
	return self.outbound
}

func (self *stackDefinition) StackState() *StackState {
	return self.stackState
}
