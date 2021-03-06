package types

type ApplicationsManager struct {
	RegisterApp               func(appName string, app AppRegistrationData)
	RunApp                    func(name string, desktopManager DesktopManager, input *interface{}) (appInstance *AppInstance, error error)
	CancelApp                 func(id int64)
	RunningApps               *[]AppInstance
	RegisteredApps            *map[string]AppRegistrationData
	ListenToRunningAppsUpdate func(func())
}

type ApplicationsManagerDependencies struct {
	PortManager PortManager
	Logger      Logger
}

type AppInstance struct {
	AppData RunningAppData
	Id      int64
	Cancel  func()
	Offline func()
	Online  func()
}
