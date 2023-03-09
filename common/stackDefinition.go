package common

type stackDefinition struct {
	name       string
	onInbound  BoundResult
	onOutbound BoundResult
	stackState IStackState
}

func NewStackDefinition(
	name string,
	inbound BoundResult,
	outbound BoundResult,
	stackState IStackState,
) (IStackDefinition, error) {
	return &stackDefinition{
		name:       name,
		onInbound:  inbound,
		onOutbound: outbound,
		stackState: stackState,
	}, nil
}

func (self *stackDefinition) Name() string {
	return self.name
}

func (self *stackDefinition) OnInbound() BoundResult {
	return self.onInbound
}

func (self *stackDefinition) OnOutbound() BoundResult {
	return self.onOutbound
}

func (self *stackDefinition) StackState() IStackState {
	return self.stackState
}
