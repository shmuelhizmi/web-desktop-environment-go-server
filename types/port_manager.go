package types

type PortManagerDependencies struct {
	Logger          Logger
	SettingsManager SettingsManager
}

type PortManager struct {
	GetDesktopPort func() (error, int32)
	GetAppPort     func() (error, int32)
}
