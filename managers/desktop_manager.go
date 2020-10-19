package managers

import (
	"github.com/shmuelhizmi/web-desktop-environment-go-server/types"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/utils"
)

func CreateDesktopManager() types.DesktopManager {
	rootLogger := utils.CreateRootLogger()
	portManager := CreatePortManager(types.PortMangerDependencies{Logger: rootLogger})
	settingsManager := CreateSettingsManager(types.SettingsManagerDependencies{Logger: rootLogger})
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
