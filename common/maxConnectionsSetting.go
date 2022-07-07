package common

type maxConnectionsSetting struct {
	maxConnections int
}

func (self *maxConnectionsSetting) ApplyNetManagerSettings(settings *NetManagerSettings) error {
	settings.MaxConnections = self.maxConnections
	return nil
}

func MaxConnectionsSetting(maxConnections int) *maxConnectionsSetting {
	return &maxConnectionsSetting{maxConnections: maxConnections}
}
