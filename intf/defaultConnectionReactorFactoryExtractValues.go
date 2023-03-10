package intf

type defaultConnectionReactorFactoryExtractValues struct {
}

func NewDefaultConnectionReactorFactoryExtractValues() IConnectionReactorFactoryExtractValues {
	return &defaultConnectionReactorFactoryExtractValues{}
}

func (self *defaultConnectionReactorFactoryExtractValues) Values(_ map[string]interface{}) (map[string]interface{}, error) {
	return make(map[string]interface{}), nil
}
