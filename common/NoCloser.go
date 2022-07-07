package common

type NoCloser struct {
}

func NewNoCloser() *NoCloser {
	return &NoCloser{}
}

func (n NoCloser) Close() error {
	return nil
}
