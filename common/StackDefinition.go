package common

type StackDefinition struct {
	name       string
	Inbound    IBoundResult
	Outbound   IBoundResult
	StackState *StackState
}

func NewStackDefinition(
	name string,
	inbound IBoundResult,
	outbound IBoundResult,
	stackState *StackState,
) (*StackDefinition, error) {
	return &StackDefinition{
		name:       name,
		Inbound:    inbound,
		Outbound:   outbound,
		StackState: stackState,
	}, nil
}

func (self *StackDefinition) GetName() string {
	return self.name
}

func (self *StackDefinition) GetInbound() IBoundResult {
	return self.Inbound
}

func (self *StackDefinition) GetOutbound() IBoundResult {
	return self.Outbound
}

func (self *StackDefinition) GetStackState() *StackState {
	return self.StackState
}
