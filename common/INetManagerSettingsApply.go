package common

type INetManagerSettingsApply interface {
	ApplyNetManagerSettings(settings *NetManagerSettings) error
}
