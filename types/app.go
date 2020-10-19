package types

import (
	react_fullstack_go_server "github.com/shmuelhizmi/react-fullstack-go-server"
)

type RunningAppData struct {
	AppPort    int32
	Name       string
	Icon       Icon
	NativeIcon NativeIcon
}

type AppRegistrationData struct {
	Icon         Icon
	NativeIcon   NativeIcon
	Name         string
	DefaultInput interface{}
	Description  string
	App          func(desktopManager DesktopManager, input interface{}) (component react_fullstack_go_server.Component)
}

type OpenApp struct {
	Name       string     `json:"name"`
	Port       int32      `json:"port"`
	Icon       Icon       `json:"icon"`
	NativeIcon NativeIcon `json:"nativeIcon"`
	Id         int64      `json:"id"`
}

type App struct {
	Name        string     `json:"name"`
	Icon        Icon       `json:"icon"`
	NativeIcon  NativeIcon `json:"nativeIcon"`
	Flow        string     `json:"flow"`
	Description string     `json:"description"`
}
