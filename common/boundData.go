package common

type boundData struct {
	name  string
	bound BoundResult
}

func (self *boundData) Name() string {
	return self.name
}

func (self *boundData) Bound() BoundResult {
	return self.bound
}

type IBoundData interface {
	Name() string
	Bound() BoundResult
}

func NewBoundData(name string, outbound BoundResult) IBoundData {
	return &boundData{
		name:  name,
		bound: outbound,
	}
}
