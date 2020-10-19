package desktop

import "github.com/shmuelhizmi/web-desktop-environment-go-server/types"

func CreateWindow(input types.CreateWindowParameters) types.CreateWindowReturn {
	isCanceled := false
	windowState := input.State
	windowView := input.View(0, "Window", nil)
	windowView.Params["name"] = input.Name
	windowView.Params["icon"] = input.Icon
	windowView.Params["title"] = input.Title
	windowView.Params["window"] = windowState
	settings := input.DesktopManager.SettingsManager.Settings()
	updateWindowParamsFromSettingsManager := func() {
		if !isCanceled {
			windowView.Params["background"] = settings.Desktop.Background
		}
	}
	updateWindowParamsFromSettingsManager()
	return types.CreateWindowReturn{
		UpdateTitle: func(newTitle string) {
			if !isCanceled {
				windowView.Params["title"] = newTitle
				windowView.Update()
			}
		},
		Cancel: nil,
	}
}
