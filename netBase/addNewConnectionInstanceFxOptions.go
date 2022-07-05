package netBase

import "go.uber.org/fx"

type addNewConnectionInstanceFxOptions struct {
	options []fx.Option
}

func (self *addNewConnectionInstanceFxOptions) apply(settings *newConnectionInstanceSettings) {
	settings.options = append(settings.options, self.options...)
}

func NewAddFxOptions(options ...fx.Option) INewConnectionInstanceSettingsApply {
	return &addNewConnectionInstanceFxOptions{
		options: options,
	}
}
