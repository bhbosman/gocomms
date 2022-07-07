package netBase

type addNewConnectionInstanceSettings struct {
	settings []INewConnectionInstanceSettingsApply
}

func (self *addNewConnectionInstanceSettings) apply(settings *newConnectionInstanceSettings) {
	for _, setting := range self.settings {
		setting.apply(settings)
	}
}

func newAddSettings(settings ...INewConnectionInstanceSettingsApply) INewConnectionInstanceSettingsApply {
	return &addNewConnectionInstanceSettings{
		settings: settings,
	}
}
