package managers

import (
	"github.com/fatih/color"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/types"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/utils"
)

func CreateDesktopManager() types.DesktopManager {
	rootLogger := utils.CreateRootLogger().Mount("desktop manager", color.FgGreen)
	settingsManager := CreateSettingsManager(types.SettingsManagerDependencies{Logger: rootLogger})
	portManager := CreatePortManager(types.PortManagerDependencies{Logger: rootLogger, SettingsManager: settingsManager})
	applicationsManager := CreateApplicationsManager(types.ApplicationsManagerDependencies{
		PortManager: portManager,
		Logger:      rootLogger,
	})
	downloadManager := CreateDownloadManager(types.DownloadManagerDependencies{
		Logger:      rootLogger,
		PortManager: portManager,
	})
	vmsManager := CreateVMSManager(types.VMSManagerDependencies{
		Logger:              rootLogger,
		PortManager:         portManager,
		SettingsManager:     settingsManager,
		ApplicationsManager: applicationsManager,
	})
	return types.DesktopManager{
		PortManager:         portManager,
		SettingsManager:     settingsManager,
		MountLogger:         rootLogger.Mount,
		ApplicationsManager: applicationsManager,
		DownloadManager:     downloadManager,
		VMSManager:          vmsManager,
	}
}
