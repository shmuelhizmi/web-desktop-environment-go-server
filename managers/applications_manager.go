package managers

import (
	"errors"
	"github.com/fatih/color"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	react_fullstack_go_server "github.com/shmuelhizmi/react-fullstack-go-server"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/types"
)

func CreateApplicationsManager(dependencies types.ApplicationsManagerDependencies) (appManager types.ApplicationsManager) {
	logger := dependencies.Logger.Mount("applications manager", color.FgMagenta)
	var runningApps []types.AppInstance
	var runningAppsUpdateListener []func()
	var appIndex int64 = 0
	callAppRunningAppsUpdateListeners := func() {
		for _, listener := range runningAppsUpdateListener {
			listener()
		}
	}
	registeredApps := make(map[string]types.AppRegistrationData)

	return types.ApplicationsManager{
		RegisterApp: func(appName string, app types.AppRegistrationData) {
			registeredApps[appName] = app
		},
		RunApp: func(name string, desktopManager types.DesktopManager, input *interface{}) (appInstance *types.AppInstance, error error) {
			app, ok := registeredApps[name]
			if !ok {
				return nil, errors.New("trying to run a non existing app")
			}
			appInput := app.DefaultInput
			if input != nil {
				appInput = *input
			}
			appId := appIndex
			appIndex++
			appMountPath := dependencies.NetworkManager.GetAppPath(name)
			appDesktopManager := desktopManager
			appDesktopManager.MountLogger = logger.Mount
			appComponent := app.App(appDesktopManager, appId, appInput)
			appServer := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())
			runningAppComponentInstance := react_fullstack_go_server.App(appServer, appComponent)
			removeAppFromRunningApps := func() {
				for appIndex, currentApp := range runningApps {
					if currentApp.Id == appId {
						runningApps = append(runningApps[:appIndex], runningApps[appIndex+1:]...)
						callAppRunningAppsUpdateListeners()
					}
				}
			}
			runningApp := types.AppInstance{
				AppData: types.RunningAppData{
					AppPath:    appMountPath,
					Name:       app.Name,
					Icon:       app.Icon,
					NativeIcon: app.NativeIcon,
				},
				Cancel: func() {
					runningAppComponentInstance.Cancel()
					removeAppFromRunningApps()
				},
				Offline: func() {
					runningAppComponentInstance.Stop()
				},
				Online: func() {
					runningAppComponentInstance.Continue()
				},
				Id: appId,
			}
			runningApps = append(runningApps, runningApp)
			callAppRunningAppsUpdateListeners()
			serveMux := dependencies.NetworkManager.Server
			serveMux.Handle("/socket.io" + appMountPath, appServer)
			logger.Info("starting app " + app.Name + "on " + appMountPath)
			return appInstance, nil
		},
		CancelApp: func(id int64) {
			appIndex := -1
			for currentAppIndex, app := range runningApps {
				if app.Id == id {
					appIndex = currentAppIndex
				}
			}
			if appIndex != -1 {
				runningApps[appIndex].Cancel()
			}
		},
		RunningApps:    &runningApps,
		RegisteredApps: &registeredApps,
		ListenToRunningAppsUpdate: func(listener func()) {
			runningAppsUpdateListener = append(runningAppsUpdateListener, listener)
		},
	}
}
