package managers

import (
	"github.com/fatih/color"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/utils"
)

type DesktopManager struct {
	PortManager PortManager
	SettingsManager SettingsManger
	MountLogger func(levelName string, levelColor color.Attribute) utils.Logger
}

func CreateDesktopManager() DesktopManager {
	rootLogger := utils.CreateRootLogger()
	portManager := CreatePortManager(PortMangerDependencies{logger: rootLogger})
	settingsManager := CreateSettingsManager(SettingsManagerDependencies{logger: rootLogger})
	return DesktopManager{
		PortManager: portManager,
		SettingsManager: settingsManager,
		MountLogger: rootLogger.Mount,
	}
}