package netListener

type netListenManagerSettings struct {
	userContext interface{}
	maxConnections int
}

type ListenAppSettingsApply interface {
	apply(settings *netListenManagerSettings)
}

type userContextValue struct {
	userContext interface{}
}

func (self userContextValue) apply(settings *netListenManagerSettings) {
	settings.userContext = self.userContext
}

func UserContextValue(userContext interface{}) *userContextValue {
	return &userContextValue{userContext: userContext}
}
type maxConnectionsSetting struct {
	maxConnections int
}

func MaxConnectionsSetting(maxConnections int) *maxConnectionsSetting {
	return &maxConnectionsSetting{maxConnections: maxConnections}
}

func (self maxConnectionsSetting) apply(settings *netListenManagerSettings) {
	settings.maxConnections = self.maxConnections
}
