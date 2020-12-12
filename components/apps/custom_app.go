package apps

import (
	react_fullstack_go_server "github.com/shmuelhizmi/react-fullstack-go-server"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/components/desktop"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/types"
)

func CreateCustomApp(customAppRegistrationData types.CustomAppRegistrationData) types.AppRegistrationData {
	return types.AppRegistrationData{
		Icon:         customAppRegistrationData.Icon,
		NativeIcon:   customAppRegistrationData.NativeIcon,
		Name:         customAppRegistrationData.Name,
		DefaultInput: nil,
		Description:  customAppRegistrationData.Description,
		App: func(desktopManager types.DesktopManager, appId int64, input interface{}) (component react_fullstack_go_server.Component) {
			window := desktop.CreateWindow(types.CreateWindowParameters{
				Name:           customAppRegistrationData.Window.Name,
				Title:          customAppRegistrationData.Window.Title,
				Icon:           customAppRegistrationData.Window.Icon,
				State:          customAppRegistrationData.Window.State,
				DesktopManager: desktopManager,
				App: CustomApp(desktopManager, types.CustomAppInput{
					Script:      customAppRegistrationData.Script,
					Permissions: customAppRegistrationData.Permissions,
				}),
				AppId: appId,
			})
			return window.Component
		},
	}
}

func CustomApp(desktopManager types.DesktopManager, input types.CustomAppInput) (component react_fullstack_go_server.Component) {
	return func(params *react_fullstack_go_server.ComponentParams) {
		desktopManager.VMSManager.RunVMApp(input.Script, input.Permissions, params)
		<-params.Cancel
	}
}
