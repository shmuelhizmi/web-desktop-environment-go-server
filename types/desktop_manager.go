package types

import (
	"github.com/fatih/color"
)

type DesktopManager struct {
	PortManager         PortManager
	SettingsManager     SettingsManger
	ApplicationsManager ApplicationsManager
	MountLogger         func(levelName string, levelColor ...color.Attribute) Logger
}
