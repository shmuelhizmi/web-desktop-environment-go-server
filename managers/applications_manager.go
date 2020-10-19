package managers

import (
	"errors"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	react_fullstack_go_server "github.com/shmuelhizmi/react-fullstack-go-server"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/types"
	"log"
	"net/http"
	"strconv"
)

func CreateApplicationsManager(dependencies types.ApplicationsManagerDependencies) (appManager types.ApplicationsManager) {
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
			if *input != nil {
				appInput = *input
			}
			getAppPortError, appPort := dependencies.PortManager.GetAppPort()
			if getAppPortError != nil {
				return nil, getAppPortError
			}
			appComponent := app.App(desktopManager, appInput)
			appServer := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())
			runningAppComponentInstance := react_fullstack_go_server.App(appServer, appComponent)
			appId := appIndex
			appIndex++
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
					AppPort:    appPort,
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
			serveMux := http.NewServeMux()
			serveMux.Handle("/socket.io/", appServer)
			log.Panic(http.ListenAndServe(":"+strconv.FormatInt(int64(appPort), 10), serveMux))
			return appInstance, nil
		},
		CancelApp: func(id int64) {
			for appIndex, app := range runningApps {
				if app.Id == id {
					app.Cancel()
					runningApps = append(runningApps[:appIndex], runningApps[appIndex+1:]...)
					callAppRunningAppsUpdateListeners()
				}
			}
		},
		RunningApps:    &runningApps,
		RegisteredApps: &registeredApps,
		ListenToRunningAppsUpdate: func(listener func()) {
			runningAppsUpdateListener = append(runningAppsUpdateListener, listener)
		},
	}
}
