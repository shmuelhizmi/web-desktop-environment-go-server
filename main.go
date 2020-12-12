package main

import (
	"github.com/fatih/color"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/shmuelhizmi/react-fullstack-go-server"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/components/apps"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/components/desktop"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/managers"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/types"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/utils"
	"log"
	"net/http"
)

func main() {

	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	desktopManager := managers.CreateDesktopManager()

	desktopManager.SettingsManager.Initialize()
	_ = desktopManager.DownloadManager.Initialize()

	customApp := apps.CreateCustomApp(types.CustomAppRegistrationData{
		Icon: types.Icon{
			Icon:     "FcFolder",
			IconType: "icon",
		},
		NativeIcon: types.NativeIcon{
			Icon:     "folder-multiple",
			IconType: "MaterialCommunityIcons",
		},
		Name:        "Custom App",
		Description: "custom app built with tengo",
		Window: types.CustomAppWindow{
			Name:  "custom app window",
			Title: "custom app",
			Icon: types.Icon{
				Icon:     "FcFolder",
				IconType: "icon",
			},
			State: types.WindowState{
				Width:     780,
				MaxWidth:  1200,
				MinWidth:  600,
				Height:    600,
				MaxHeight: 800,
				MinHeight: 450,
				Position: types.Position{
					X: 150,
					Y: 150,
				},
			},
		},
		Permissions: []types.VMPermission{types.OSPermission, types.FmtPermission},
		Script: `
		app := import("app")
		logger := import("fmt")
		explorer := app.view(0, "Explorer")
		explorer.updateParam("platformPathSeparator", "/")
		explorer.updateParam("type", "explore")
		explorer.updateParam("files", [{ isFolder: false, name: "hello.txt" }])
		explorer.updateParam("currentPath", "/home/user/")
		explorer.handleFunction("onChangeCurrentPath",
								func (newPath) {
									logger.println(newPath + " is the new path")
									explorer.updateParam("currentPath", newPath)
									explorer.update()
									return true
								}
							);
		logger.println("starting view")
		explorer.start()
`,
	})

	desktopManager.ApplicationsManager.RegisterApp(customApp.Name, customApp)
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
