package netDial

type dialAppSettings struct {
	userContext interface{}
	canDial     []ICanDial
	maxConnections int
}

type DialAppSettingsApply interface {
	apply(settings *dialAppSettings)
}

type userContextValue struct {
	userContext interface{}
}

func (self userContextValue) apply(settings *dialAppSettings) {
	settings.userContext = self.userContext
}

func UserContextValue(userContext interface{}) *userContextValue {
	return &userContextValue{userContext: userContext}
}

type canDialSetting struct {
	canDial []ICanDial
}

func CanDial(canDial ...ICanDial) *canDialSetting {
	return &canDialSetting{canDial: canDial}
}

func (self canDialSetting) apply(settings *dialAppSettings) {
	for _, cd := range self.canDial{
		settings.canDial = append(settings.canDial, cd)
	}
}



type maxConnectionsSetting struct {
	maxConnections int
}

func MaxConnectionsSetting(maxConnections int) *maxConnectionsSetting {
	return &maxConnectionsSetting{maxConnections: maxConnections}
}

func (self maxConnectionsSetting) apply(settings *dialAppSettings) {
	settings.maxConnections = self.maxConnections
}
