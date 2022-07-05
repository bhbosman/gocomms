package common

import (
	"github.com/bhbosman/goCommsDefinitions"
	"go.uber.org/fx"
)

type overrideOnCreateConnectionFactory struct {
	onCreateConnectionFactory func() (goCommsDefinitions.IOnCreateConnection, error)
}

func (self *overrideOnCreateConnectionFactory) ApplyNetManagerSettings(settings *NetManagerSettings) error {
	settings.SetOnCreateConnection(self.onCreateConnectionFactory)
	return nil
}

func NewOverrideOnCreateConnectionFactory(listenerFactory func() (goCommsDefinitions.IOnCreateConnection, error)) *overrideOnCreateConnectionFactory {
	return &overrideOnCreateConnectionFactory{
		onCreateConnectionFactory: listenerFactory,
	}
}

type overrideCreateConnectionReactor struct {
	provide fx.Option
}

func (self *overrideCreateConnectionReactor) ApplyNetManagerSettings(settings *NetManagerSettings) error {
	settings.FxOptionsForConnectionInstance = append(settings.FxOptionsForConnectionInstance, self.provide)
	return nil
}

func NewOverrideCreateConnectionReactor(provide fx.Option) *overrideCreateConnectionReactor {
	return &overrideCreateConnectionReactor{provide: provide}
}
