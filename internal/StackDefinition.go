package internal

import "github.com/google/uuid"

type BoundResult func(params InOutBoundParams) (IStackBoundDefinition, error)

type IBoundResult interface {
	GetBoundResult() (BoundResult, error)
}

type boundResultImpl struct {
	boundResult BoundResult
}

func (self *boundResultImpl) GetBoundResult() (BoundResult, error) {
	return self.boundResult, nil
}

func NewBoundResultImpl(boundResult BoundResult) *boundResultImpl {
	return &boundResultImpl{boundResult: boundResult}
}

type IStackDefinition interface {
	GetId() uuid.UUID
	GetName() string
	GetInbound() IBoundResult
	GetOutbound() IBoundResult
	GetStackState() *StackState
}

type StackDefinition struct {
	IId        uuid.UUID
	Name       string
	Inbound    IBoundResult
	Outbound   IBoundResult
	StackState *StackState
}

func (self *StackDefinition) GetId() uuid.UUID {
	return self.IId
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
