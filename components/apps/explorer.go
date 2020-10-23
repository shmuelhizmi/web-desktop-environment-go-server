package apps

import (
	"encoding/json"
	"github.com/fatih/color"
	react_fullstack_go_server "github.com/shmuelhizmi/react-fullstack-go-server"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/components/desktop"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/types"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/utils"
	"os"
)

func GetExplorerAppInfo() types.RegisteredApp {
	return types.RegisteredApp{
		Name: "explorer",
		Icon: types.Icon{
			Icon:     "FcFolder",
			IconType: "icon",
		},
		NativeIcon: types.NativeIcon{
			Icon:     "folder-multiple",
			IconType: "MaterialCommunityIcons",
		},
		RegisteredName: "explorer",
		Description:    "a file explorer",
	}
}

func CreateExplorerApp() types.AppRegistrationData {
	home, _ := os.UserHomeDir()
	info := GetExplorerAppInfo()
	return types.AppRegistrationData{
		Icon:       info.Icon,
		NativeIcon: info.NativeIcon,
		Name:       info.Name,
		DefaultInput: types.ExplorerInput{
			Path:         home,
			IsCurrentApp: true,
			Type:         "explore",
		},
		Description: info.Description,
		App: func(desktopManager types.DesktopManager, appId int64, input interface{}) (component react_fullstack_go_server.Component) {
			var explorerInput types.ExplorerInput = input.(types.ExplorerInput)
			window := desktop.CreateWindow(types.CreateWindowParameters{
				Name:  info.Name,
				Title: info.Name,
				Icon:  info.Icon,
				State: types.WindowState{
					Width:     720,
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
				DesktopManager: desktopManager,
				App:            ExplorerApp(desktopManager, explorerInput),
				AppId:          appId,
			})
			return window.Component
		},
	}
}

func ExplorerApp(desktopManager types.DesktopManager, input types.ExplorerInput) (component react_fullstack_go_server.Component) {
	return func(params *react_fullstack_go_server.ComponentParams) {
		logger := desktopManager.MountLogger("explorer", color.BgBlue, color.FgBlack)
		logger.Info("running explorer")
		currentPath := input.Path
		explorerView := params.View(0, "Explorer", nil)
		explorerView.Params["platformPathSeparator"] = string(os.PathSeparator)
		explorerView.Params["type"] = input.Type
		updateExplorerWithCurrentPath := func() {
			explorerView.Params["currentPath"] = currentPath
			explorerView.Params["files"] = utils.ListFilesInDir(currentPath)
		}
		updateExplorerWithCurrentPath()
		explorerView.On("onChangeCurrentPath", func(props [][]byte) interface{} {
			var newPath string
			json.Unmarshal(props[0], &newPath)
			currentPath = newPath
			updateExplorerWithCurrentPath()
			explorerView.Update()
			return nil
		})
		explorerView.On("onCopy", func(props [][]byte) interface{} {
			var arguments struct {
				NewPath      string `json:"newPath"`
				OriginalPath string `json:"originalPath"`
			}
			json.Unmarshal(props[0], &arguments)
			copyError := utils.Copy(arguments.OriginalPath, arguments.NewPath)
			if copyError == nil {
				updateExplorerWithCurrentPath()
				explorerView.Update()
			} else {
				logger.Error("fail to copy " + arguments.OriginalPath + " to " + arguments.NewPath)
				logger.Error(copyError.Error())
			}
			return nil
		})
		explorerView.On("onRequestDownloadLink", func(props [][]byte) interface{} {
			return struct {
				Path string `json:"path"`
				Port int32  `json:"port"`
			}{
				Path: "/fwafafwawaf",
				Port: 0,
			}
		})
		explorerView.Start()
		<-params.Cancel
	}
}
