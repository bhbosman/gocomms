package common

type IBoundResult interface {
	GetBoundResult() (BoundResult, error)
}
