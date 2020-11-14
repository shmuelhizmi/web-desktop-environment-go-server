package desktop

import (
	"github.com/fatih/color"
	react_fullstack_go_server "github.com/shmuelhizmi/react-fullstack-go-server"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/types"
)

func CreateDesktop(desktopManager types.DesktopManager) react_fullstack_go_server.Component {
	return func(params *react_fullstack_go_server.ComponentParams) {
		settings := desktopManager.SettingsManager.Settings()
		desktopLogger := desktopManager.MountLogger("desktop", color.FgHiCyan)
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
			apps := make([]types.RegisteredApp, 0, len(*registeredApps))
			for appName, app := range *registeredApps {
				apps = append(apps, types.RegisteredApp{
					Name:           app.Name,
					Icon:           app.Icon,
					NativeIcon:     app.NativeIcon,
					RegisteredName: appName,
					Description:    app.Description,
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
		desktop.On("onLaunchApp", func(app struct {
			App    string      `json:"flow"`
			Params interface{} `json:"params"` // deprecated
		}) {
			desktopLogger.Info("launching app " + app.App)
			desktopManager.ApplicationsManager.RunApp(app.App, desktopManager, nil)
		})
		desktop.On("onCloseApp", func(appId int64) interface{} {
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
		<-params.Cancel
	}
}
