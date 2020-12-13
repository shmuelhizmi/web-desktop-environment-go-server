package managers

import (
	"github.com/fatih/color"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/types"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/utils"
)

func CreateDesktopManager() types.DesktopManager {
	rootLogger := utils.CreateRootLogger().Mount("desktop manager", color.FgGreen)
	settingsManager := CreateSettingsManager(types.SettingsManagerDependencies{Logger: rootLogger})
	networkManger := CreateNetworkManager(types.NetworkManagerDependencies{
		SettingsManager: settingsManager,
		Logger:          rootLogger,
	})
	applicationsManager := CreateApplicationsManager(types.ApplicationsManagerDependencies{
		Logger:         rootLogger,
		NetworkManager: networkManger,
	})
	downloadManager := CreateDownloadManager(types.DownloadManagerDependencies{
		Logger:         rootLogger,
		NetworkManager: networkManger,
	})
	return types.DesktopManager{
		SettingsManager:     settingsManager,
		MountLogger:         rootLogger.Mount,
		ApplicationsManager: applicationsManager,
		DownloadManager:     downloadManager,
		NetworkManager:      networkManger,
	}
}
