package apps

import (
	"github.com/fatih/color"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	react_fullstack_go_server "github.com/shmuelhizmi/react-fullstack-go-server"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/components/desktop"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/types"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/utils"
	"os"
)

func GetTerminalAppInfo() types.RegisteredApp {
	return types.RegisteredApp{
		Name: "Terminal",
		Icon: types.Icon{
			Icon:     "FcCommandLine",
			IconType: "icon",
		},
		NativeIcon: types.NativeIcon{
			Icon:     "console",
			IconType: "MaterialCommunityIcons",
		},
		RegisteredName: "terminal",
		Description:    "a terminal window",
	}
}

func CreateTerminalApp() types.AppRegistrationData {
	info := GetTerminalAppInfo()
	home, getHomedirError := os.UserHomeDir()
	if getHomedirError != nil {
		home = "/home/"
	}
	return types.AppRegistrationData{
		Icon:       info.Icon,
		NativeIcon: info.NativeIcon,
		Name:       info.Name,
		DefaultInput: types.TerminalInput{
			Process:  "/bin/bash",
			Argument: make([]string, 0),
			Location: home,
		},
		Description: info.Description,
		App: func(desktopManager types.DesktopManager, appId int64, input interface{}) (component react_fullstack_go_server.Component) {
			window := desktop.CreateWindow(types.CreateWindowParameters{
				Name:  info.Name,
				Title: info.Name,
				Icon:  info.Icon,
				State: types.WindowState{
					Width:     1000,
					MaxWidth:  1200,
					MinWidth:  350,
					Height:    400,
					MaxHeight: 900,
					MinHeight: 200,
					Position: types.Position{
						X: 50,
						Y: 50,
					},
				},
				DesktopManager: desktopManager,
				App:            TerminalApp(desktopManager, input.(types.TerminalInput), appId),
				AppId:          appId,
			})
			return window.Component
		},
	}
}

func TerminalApp(desktopManager types.DesktopManager, input types.TerminalInput, appId int64) (component react_fullstack_go_server.Component) {
	return func(params *react_fullstack_go_server.ComponentParams) {
		logger := desktopManager.MountLogger("terminal", color.BgBlue, color.FgBlack)
		logger.Info("running terminal")
		terminalView := params.View(0, types.TerminalViewName, nil)
		terminalInstance, createTerminalError := utils.CreatePtyInstance(input.Process, input.Argument, input.Location, logger)
		if createTerminalError != nil {
			desktopManager.ApplicationsManager.CancelApp(appId)
			logger.Error("fail to start terminal instance - closing app")
		} else {
			socketServerPath := desktopManager.NetworkManager.GetServicePath("terminal")
			server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())
			serveMux := desktopManager.NetworkManager.Server
			serveMux.Handle("/socket.io" + socketServerPath, server)
			history := ""
			server.On(gosocketio.OnConnection, func(sender *gosocketio.Channel) {
				sender.Emit("output", history)
			})
			server.On("input", func(sender *gosocketio.Channel, input string) {
				terminalInstance.Write(input)
			})
			terminalInstance.ListenToOutput(func(newOutput string) {
				server.BroadcastToAll("output", newOutput)
				history += newOutput
			})
			terminalView.Params["path"] = socketServerPath

		}
		terminalView.Start()
		<-params.Cancel
		if createTerminalError != nil {
			terminalInstance.Exit()
		}
	}
}
