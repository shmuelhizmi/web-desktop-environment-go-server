package types

import (
	"github.com/fatih/color"
)

type DesktopManager struct {
	PortManager         PortManager
	SettingsManager     SettingsManager
	ApplicationsManager ApplicationsManager
	DownloadManager     DownloadManager
	VMSManager          VMSManager
	MountLogger         func(levelName string, levelColor ...color.Attribute) Logger
}
