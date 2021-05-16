package internal

type InOutBoundParams struct {
	Index int
}

func NewInOutBoundParams(index int) InOutBoundParams {
	return InOutBoundParams{Index: index}
}
