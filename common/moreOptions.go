package common

import "go.uber.org/fx"

type moreOptions struct {
	options []fx.Option
}

func (self *moreOptions) ApplyNetManagerSettings(settings *NetManagerSettings) error {
	settings.AddMoreOptions(self.options)
	return nil
}

func MoreOptions(options ...fx.Option) *moreOptions {
	return &moreOptions{
		options: options,
	}
}
