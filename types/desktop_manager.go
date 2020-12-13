package types

import (
	"github.com/fatih/color"
)

type DesktopManager struct {
	SettingsManager     SettingsManager
	NetworkManager      NetworkManager
	ApplicationsManager ApplicationsManager
	DownloadManager     DownloadManager
	MountLogger         func(levelName string, levelColor ...color.Attribute) Logger
}
