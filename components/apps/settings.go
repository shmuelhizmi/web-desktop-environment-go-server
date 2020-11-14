package apps

import (
	"errors"
	"github.com/fatih/color"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	react_fullstack_go_server "github.com/shmuelhizmi/react-fullstack-go-server"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/components/desktop"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/types"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/utils"
)

func GetSettingsAppInfo() types.RegisteredApp {
	return types.RegisteredApp{
		Name: "Settings",
		Icon: types.Icon{
			Icon:     "FcSettings",
			IconType: "icon",
		},
		NativeIcon: types.NativeIcon{
			Icon:     "ios-settings",
			IconType: "Ionicons",
		},
		RegisteredName: "explorer",
		Description:    "a file explorer",
	}
}

func CreateSettingsApp() types.AppRegistrationData {
	info := GetSettingsAppInfo()
	return types.AppRegistrationData{
		Icon:         info.Icon,
		NativeIcon:   info.NativeIcon,
		Name:         info.Name,
		DefaultInput: struct{}{},
		Description:  info.Description,
		App: func(desktopManager types.DesktopManager, appId int64, input interface{}) (component react_fullstack_go_server.Component) {
			window := desktop.CreateWindow(types.CreateWindowParameters{
				Name:  info.Name,
				Title: info.Name,
				Icon:  info.Icon,
				State: types.WindowState{
					Width:     920,
					MaxWidth:  1200,
					MinWidth:  500,
					Height:    600,
					MaxHeight: 900,
					MinHeight: 600,
					Position: types.Position{
						X: 50,
						Y: 50,
					},
				},
				DesktopManager: desktopManager,
				App:            SettingsApp(desktopManager, input.(struct{})),
				AppId:          appId,
			})
			return window.Component
		},
	}
}

func SettingsApp(desktopManager types.DesktopManager, input struct{}) (component react_fullstack_go_server.Component) {
	return func(params *react_fullstack_go_server.ComponentParams) {
		logger := desktopManager.MountLogger("settings", color.BgBlue, color.FgBlack)
		logger.Info("running settings")
		settingsView := params.View(0, types.SettingsViewName, nil)
		updateSettingsViewWithSettings := func() {
			settings := desktopManager.SettingsManager.Settings()
			settingsView.Params["settings"] = *settings
		}
		updateSettingsViewWithSystemInfo := func() error {

			hostStat, getHostError := host.Info()
			cpuStat, getCpuError := cpu.Info()
			vmStat, getMemoryError := mem.VirtualMemory()
			if getHostError != nil || getCpuError != nil || getMemoryError != nil || len(cpuStat) == 0 {
				logger.Error("fail to fetch system information for settings")
				return errors.New("fail to get some information")
			}
			settingsView.Params["systemInfo"] = types.SystemInformation{
				CPU: types.CPUInformation{
					BrandName:     cpuStat[0].ModelName,
					Speed:         utils.FloatToString(cpuStat[0].Mhz / 1000),
					PhysicalCores: utils.Int32ToString(int32(len(cpuStat))),
					Cores:         utils.Int32ToString(int32(len(cpuStat)) * cpuStat[0].Cores),
				},
				Ram: types.RamInformation{
					Total: utils.UInt64ToString(vmStat.Total/1024/1024/1024) + "GB",
					Free:  utils.UInt64ToString(vmStat.Available/1024/1024/1024) + "GB",
				},
				Disks: nil,
				OS: types.OSInformation{
					Platform: hostStat.Platform,
					Kernel:   hostStat.KernelArch,
					Hostname: hostStat.Hostname,
				},
			}
			settingsView.Update()
			return nil
		}
		updateSettingsViewWithSettings()
		settingsView.On("onReload", func() {
			updateSettingsViewWithSettings()
			settingsView.Update()
			go updateSettingsViewWithSystemInfo()
		})
		settingsView.On("setSettings", func(newSettings types.SettingsObject) {
			desktopManager.SettingsManager.SetSettings(newSettings)
			updateSettingsViewWithSettings()
		})
		settingsView.Start()
		go updateSettingsViewWithSystemInfo()
		<-params.Cancel
	}
}
