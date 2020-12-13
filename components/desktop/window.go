package desktop

import (
	react_fullstack_go_server "github.com/shmuelhizmi/react-fullstack-go-server"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/types"
)

func CreateWindow(input types.CreateWindowParameters) types.CreateWindowReturn {
	updateTitle := func(newTitle string) {
		//mock
	}

	return types.CreateWindowReturn{
		UpdateTitle: &updateTitle,
		Component: func(params *react_fullstack_go_server.ComponentParams) {
			isCanceled := false
			windowState := input.State
			themeProvider := params.View(0, "ThemeProvider", nil)
			settings := input.DesktopManager.SettingsManager.Settings()
			updateThemeProviderParamsFromSettings := func() {
				themeProvider.Params["theme"] = settings.Desktop.Theme
				themeProvider.Params["customTheme"] = settings.Desktop.CustomTheme
			}
			updateThemeProviderParamsFromSettings()
			themeProvider.Start()
			windowView := params.View(0, "Window", &themeProvider)
			windowView.Params["name"] = input.Name
			windowView.Params["icon"] = input.Icon
			windowView.Params["title"] = input.Title
			windowView.Params["window"] = windowState
			updateWindowParamsFromSettingsManager := func() {
				if !isCanceled {
					windowView.Params["background"] = settings.Desktop.Background
				}
			}
			updateWindowParamsFromSettingsManager()
			windowView.On("setWindowState", func(newWindowState struct {
				Minimized bool           `json:"minimized"`
				Position  types.Position `json:"position"`
				Size      types.Size     `json:"size"`
			}) interface{} {
				windowState.Position = newWindowState.Position
				windowState.Height = newWindowState.Size.Height
				windowState.Width = newWindowState.Size.Width
				windowState.Position = newWindowState.Position
				windowView.Params["window"] = windowState
				windowView.Update()
				return nil
			})
			input.DesktopManager.SettingsManager.ListenToNewSettings(func(_ *types.SettingsObject) {
				if !isCanceled {
					updateWindowParamsFromSettingsManager()
					updateThemeProviderParamsFromSettings()
					themeProvider.Update()
					windowView.Update()
				}
			})
			windowView.On("onClose", func() interface{} {
				input.DesktopManager.ApplicationsManager.CancelApp(input.AppId)
				return nil
			})
			updateTitle = func(newTitle string) {
				if !isCanceled {
					windowView.Params["title"] = newTitle
				}
			}
			windowView.Start()
			params.Run(input.App, windowView)
			<-params.Cancel
			isCanceled = true
		},
	}
}
