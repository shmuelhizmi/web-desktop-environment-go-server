package types

import (
	"github.com/fatih/color"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/utils"
)

type DesktopManager struct {
	PortManager         PortManager
	SettingsManager     SettingsManger
	ApplicationsManager ApplicationsManager
	MountLogger         func(levelName string, levelColor color.Attribute) utils.Logger
}
