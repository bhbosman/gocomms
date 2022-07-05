package common

import (
	"github.com/bhbosman/goCommsDefinitions"
	"go.uber.org/fx"
)

type NetManagerSettings struct {
	MaxConnections                 int
	OnCreateConnectionFactory      interface{} //func(IOnCreateConnection, err)
	OnCreateIConnectionReactor     fx.Option   //func(intf.IConnectionReactor, err)
	FxOptionsForConnectionInstance []fx.Option
	MoreOptions                    []fx.Option
}

func (self *NetManagerSettings) SetOnCreateConnection(factory func() (goCommsDefinitions.IOnCreateConnection, error)) {
	self.OnCreateConnectionFactory = factory
}

func (self *NetManagerSettings) AddFxOptionsForConnectionInstance(options []fx.Option) {
	self.FxOptionsForConnectionInstance = append(self.FxOptionsForConnectionInstance, options...)
}

func (self *NetManagerSettings) AddMoreOptions(options []fx.Option) {
	self.MoreOptions = append(self.MoreOptions, options...)
}

func (self *NetManagerSettings) Build() (func() fx.Option, error) {
	return func() fx.Option {
		var o []fx.Option
		o = append(o, self.FxOptionsForConnectionInstance...)
		return fx.Options(o...)
	}, nil
}

// NewNetManagerSettings
// ResponsibleForClientContext params is important where the same Connection Reactor might be used again and again,
// like in a SSH connection, where the same connection Reactor is used for the actual connections and for the accepted channels
func NewNetManagerSettings(
	maxConnections int,
) NetManagerSettings {
	return NetManagerSettings{
		MaxConnections:                 maxConnections,
		OnCreateConnectionFactory:      ProvideIOnCreateConnectionResource,
		FxOptionsForConnectionInstance: nil,
	}
}
