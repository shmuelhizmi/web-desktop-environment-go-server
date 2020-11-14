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
	"log"
	"net/http"
)

func main() {

	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	desktopManager := managers.CreateDesktopManager()

	desktopManager.SettingsManager.Initialize()

	desktopManager.ApplicationsManager.RegisterApp(apps.GetExplorerAppInfo().Name, apps.CreateExplorerApp())
	desktopManager.ApplicationsManager.RegisterApp(apps.GetNotepadAppInfo().Name, apps.CreateNotepadApp())
	desktopManager.ApplicationsManager.RegisterApp(apps.GetSettingsAppInfo().Name, apps.CreateSettingsApp())
	desktopManager.ApplicationsManager.RegisterApp(apps.GetTerminalAppInfo().Name, apps.CreateTerminalApp())

	mainLogger := desktopManager.MountLogger("main", color.FgHiRed)

	getDesktopPortError, desktopPort := desktopManager.PortManager.GetDesktopPort()

	if getDesktopPortError != nil {
		mainLogger.Error("we cloud not found open port to run desktop on")
		panic(getDesktopPortError)
	}

	mainLogger.Info("running app desktop on port " + utils.Int32ToString(desktopPort))

	react_fullstack_go_server.App(server, desktop.CreateDesktop(desktopManager))

	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)
	log.Panic(http.ListenAndServe(":"+utils.Int32ToString(desktopPort), serveMux))
}
