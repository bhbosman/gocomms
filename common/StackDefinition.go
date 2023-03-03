package common

type StackDefinition struct {
	Name       string
	Inbound    IBoundResult
	Outbound   IBoundResult
	StackState *StackState
}

func NewStackDefinition(
	name string,
	inbound IBoundResult,
	outbound IBoundResult,
	stackState *StackState,
) *StackDefinition {
	return &StackDefinition{
		Name:       name,
		Inbound:    inbound,
		Outbound:   outbound,
		StackState: stackState,
	}
}

func (self *StackDefinition) GetId() string {
	return self.Name
}

func (self *StackDefinition) GetName() string {
	return self.Name
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
