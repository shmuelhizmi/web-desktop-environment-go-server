package types

import (
	"github.com/fatih/color"
)

type DesktopManager struct {
	PortManager         PortManager
	SettingsManager     SettingsManager
	ApplicationsManager ApplicationsManager
	DownloadManager     DownloadManager
	MountLogger         func(levelName string, levelColor ...color.Attribute) Logger
}
