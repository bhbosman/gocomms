package common

type StackDefinition struct {
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
) (*StackDefinition, error) {
	return &StackDefinition{
		name:       name,
		inbound:    inbound,
		outbound:   outbound,
		stackState: stackState,
	}, nil
}

func (self *StackDefinition) GetName() string {
	return self.name
}

func (self *StackDefinition) GetInbound() IBoundResult {
	return self.inbound
}

func (self *StackDefinition) GetOutbound() IBoundResult {
	return self.outbound
}

func (self *StackDefinition) GetStackState() *StackState {
	return self.stackState
}
