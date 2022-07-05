package netBase

import "go.uber.org/fx"

type overrideCreateIConnectionReactor struct {
	createIConnectionReactor fx.Option
}

func (self *overrideCreateIConnectionReactor) apply(settings *newConnectionInstanceSettings) {
	settings.ProvideCreateIConnectionReactor = self.createIConnectionReactor
}

func NewOverrideCreateIConnectionReactor(createIConnectionReactor fx.Option) INewConnectionInstanceSettingsApply {
	return &overrideCreateIConnectionReactor{
		createIConnectionReactor: createIConnectionReactor,
	}
}
