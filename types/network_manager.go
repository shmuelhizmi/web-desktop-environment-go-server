package types

import "net/http"

type NetworkManager struct {
	Server         *http.ServeMux
	Initialize     func()
	GetAppPath     func(appName string) string
	GetServicePath func(service string) string
}

type NetworkManagerDependencies struct {
	SettingsManager SettingsManager
	Logger          Logger
}
