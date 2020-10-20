package managers

import (
	"github.com/shmuelhizmi/web-desktop-environment-go-server/types"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/utils"
)

func CreateDesktopManager() types.DesktopManager {
	rootLogger := utils.CreateRootLogger()
	settingsManager := CreateSettingsManager(types.SettingsManagerDependencies{Logger: rootLogger})
	portManager := CreatePortManager(types.PortMangerDependencies{Logger: rootLogger, SettingsManager: settingsManager})
	applicationsManager := CreateApplicationsManager(types.ApplicationsManagerDependencies{
		PortManager: portManager,
		Logger:      rootLogger,
	})
	return types.DesktopManager{
		PortManager:         portManager,
		SettingsManager:     settingsManager,
		MountLogger:         rootLogger.Mount,
		ApplicationsManager: applicationsManager,
	}
}
