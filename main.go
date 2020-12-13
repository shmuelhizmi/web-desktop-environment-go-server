package main

import (
	"github.com/fatih/color"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/shmuelhizmi/react-fullstack-go-server"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/components/apps"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/components/desktop"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/managers"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/utils"
	"time"
)

func main() {

	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	desktopManager := managers.CreateDesktopManager()

	desktopManager.SettingsManager.Initialize()
	desktopManager.NetworkManager.Initialize()
	_ = desktopManager.DownloadManager.Initialize()

	desktopManager.ApplicationsManager.RegisterApp(apps.GetExplorerAppInfo().Name, apps.CreateExplorerApp())
	desktopManager.ApplicationsManager.RegisterApp(apps.GetNotepadAppInfo().Name, apps.CreateNotepadApp())
	desktopManager.ApplicationsManager.RegisterApp(apps.GetSettingsAppInfo().Name, apps.CreateSettingsApp())
	desktopManager.ApplicationsManager.RegisterApp(apps.GetTerminalAppInfo().Name, apps.CreateTerminalApp())

	mainLogger := desktopManager.MountLogger("main", color.FgHiRed)

	serveMux := desktopManager.NetworkManager.Server

	mainLogger.Info(
		"running app desktop on port " +
			utils.Int32ToString(desktopManager.SettingsManager.Settings().Network.Ports.MainPort) +
			" and path /socket.io/desktop/")

	react_fullstack_go_server.App(server, desktop.CreateDesktop(desktopManager))

	serveMux.Handle("/socket.io/desktop/", server)
	for true {
		time.Sleep(time.Hour)
	}
}
