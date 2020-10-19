package desktop

import (
	"encoding/json"
	react_fullstack_go_server "github.com/shmuelhizmi/react-fullstack-go-server"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/types"
)

func CreateDesktop(desktopManager types.DesktopManager) react_fullstack_go_server.Component {
	return func(params *react_fullstack_go_server.ComponentParams) {
		settings := desktopManager.SettingsManager.Settings()
		themeProvider := params.View(0, "ThemeProvider", nil)
		updateThemeProviderParamsFromSettings := func() {
			themeProvider.Params["theme"] = settings.Desktop.Theme
			themeProvider.Params["CustomTheme"] = settings.Desktop.CustomTheme
		}
		updateThemeProviderParamsFromSettings()
		themeProvider.Start()
		desktopManager.SettingsManager.ListenToNewSettings(func(_ *types.SettingsObject) {
			updateThemeProviderParamsFromSettings()
			themeProvider.Update()
		})
		desktop := params.View(0, "Desktop", &themeProvider)
		runningApps := desktopManager.ApplicationsManager.RunningApps
		registeredApps := desktopManager.ApplicationsManager.RegisteredApps
		updateDesktopParamsFromWindowManager := func() {
			apps := make([]types.App, 0, len(*registeredApps))
			for appName, app := range *registeredApps {
				apps = append(apps, types.App{
					Name:        app.Name,
					Icon:        app.Icon,
					NativeIcon:  app.NativeIcon,
					Flow:        appName,
					Description: app.Description,
				})
			}
			openApps := make([]types.OpenApp, 0, len(*runningApps))
			for _, app := range *runningApps {
				openApps = append(openApps, types.OpenApp{
					Name:       app.AppData.Name,
					Port:       app.AppData.AppPort,
					Icon:       app.AppData.Icon,
					NativeIcon: app.AppData.NativeIcon,
					Id:         app.Id,
				})
			}
			desktop.Params["apps"] = apps
			desktop.Params["openApps"] = openApps
		}
		updateDesktopParamsFromWindowManager()
		updateDesktopParamsFromSettingsManager := func() {
			desktop.Params["background"] = settings.Desktop.Background
			desktop.Params["nativeBackground"] = settings.Desktop.NativeBackground
		}
		updateDesktopParamsFromSettingsManager()
		desktop.On("onLaunchApp", func(params [][]byte) interface{} {
			var app struct {
				App    string      `json:"flow"`
				Params interface{} `json:"params"`
			}
			json.Unmarshal(params[0], &app)
			desktopManager.ApplicationsManager.RunApp(app.App, desktopManager, &app.Params)
			return nil
		})
		desktop.On("onCloseApp", func(params [][]byte) interface{} {
			var appId int64
			json.Unmarshal(params[0], &appId)
			desktopManager.ApplicationsManager.CancelApp(appId)
			return nil
		})
		desktop.Start()
		desktopManager.SettingsManager.ListenToNewSettings(func(_ *types.SettingsObject) {
			updateDesktopParamsFromSettingsManager()
			desktop.Update()
		})
		desktopManager.ApplicationsManager.ListenToRunningAppsUpdate(func() {
			updateDesktopParamsFromWindowManager()
			desktop.Update()
		})
	}
}
