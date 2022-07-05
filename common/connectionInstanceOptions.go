package common

import (
	"go.uber.org/fx"
)

type connectionInstanceOptions struct {
	options []fx.Option
}

func (self *connectionInstanceOptions) ApplyNetManagerSettings(
	settings *NetManagerSettings,
) error {
	settings.AddFxOptionsForConnectionInstance(self.options)
	return nil
}

func NewConnectionInstanceOptions(options ...fx.Option) INetManagerSettingsApply {
	return &connectionInstanceOptions{
		options: options,
	}
}
