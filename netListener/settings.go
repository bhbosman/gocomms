package netListener

type netListenManagerSettings struct {
	userContext interface{}
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
