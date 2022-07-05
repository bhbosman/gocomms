package common

type boundResultImpl struct {
	boundResult BoundResult
}

func (self *boundResultImpl) GetBoundResult() (BoundResult, error) {
	return self.boundResult, nil
}

func NewBoundResultImpl(boundResult BoundResult) *boundResultImpl {
	return &boundResultImpl{boundResult: boundResult}
}
