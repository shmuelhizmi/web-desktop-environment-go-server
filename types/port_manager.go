package types

import "github.com/shmuelhizmi/web-desktop-environment-go-server/utils"

type PortMangerDependencies struct {
	Logger          utils.Logger
	SettingsManager SettingsManger
}

type PortManager struct {
	GetDesktopPort func() (error, int32)
	GetAppPort     func() (error, int32)
}
