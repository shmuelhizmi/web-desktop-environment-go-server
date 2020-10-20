package types

type PortMangerDependencies struct {
	Logger          Logger
	SettingsManager SettingsManger
}

type PortManager struct {
	GetDesktopPort func() (error, int32)
	GetAppPort     func() (error, int32)
}
