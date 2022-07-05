package intf

type IConnectionReactorFactoryExtractValues interface {
	Values(inputValues map[string]interface{}) (map[string]interface{}, error)
}
