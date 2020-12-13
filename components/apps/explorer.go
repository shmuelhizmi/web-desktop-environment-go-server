package apps

import (
	"github.com/fatih/color"
	react_fullstack_go_server "github.com/shmuelhizmi/react-fullstack-go-server"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/components/desktop"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/types"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/utils"
	"os"
)

func GetExplorerAppInfo() types.RegisteredApp {
	return types.RegisteredApp{
		Name: "Explorer",
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
			explorerInput := input.(types.ExplorerInput)
			window := desktop.CreateWindow(types.CreateWindowParameters{
				Name:  info.Name,
				Title: info.Name,
				Icon:  info.Icon,
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
		explorerView := params.View(0, types.ExplorerViewName, nil)
		explorerView.Params["platformPathSeparator"] = string(os.PathSeparator)
		explorerView.Params["type"] = input.Type
		updateExplorerWithCurrentPath := func() {
			explorerView.Params["currentPath"] = currentPath
			explorerView.Params["files"] = utils.FSListFilesInDir(currentPath)
		}
		updateExplorerWithCurrentPath()
		explorerView.On("onChangeCurrentPath", func(newPath string) {
			currentPath = newPath
			updateExplorerWithCurrentPath()
			explorerView.Update()
		})
		explorerView.On("onCopy", func(arguments struct {
			NewPath      string `json:"newPath"`
			OriginalPath string `json:"originalPath"`
		}) {
			copyError := utils.FSCopy(arguments.OriginalPath, arguments.NewPath)
			if copyError == nil {
				updateExplorerWithCurrentPath()
				explorerView.Update()
			} else {
				logger.Error("fail to copy " + arguments.OriginalPath + " to " + arguments.NewPath)
				logger.Error(copyError.Error())
			}
		})
		explorerView.On("onRequestDownloadLink", func(path string) interface{} {
			fileHash := desktopManager.DownloadManager.AddFile(path)
			return struct {
				Path string `json:"path"`
			}{
				Path: *desktopManager.DownloadManager.Path + fileHash,
			}
		})
		explorerView.On("onCreateFile", func(path string) {
			createFileError := utils.FSCreateEmptyFile(path)
			if createFileError == nil {
				updateExplorerWithCurrentPath()
				explorerView.Update()
			} else {
				logger.Error("fail to create empty file " + path)
				logger.Error(createFileError.Error())
			}
		})
		explorerView.On("onCreateFolder", func(path string) {
			createFolderError := utils.FSCreateFolder(path)
			if createFolderError == nil {
				updateExplorerWithCurrentPath()
				explorerView.Update()
			} else {
				logger.Error("fail to create empty folder " + path)
				logger.Error(createFolderError.Error())
			}
		})
		explorerView.On("onDelete", func(path string) {
			deleteError := utils.FSDelete(path)
			if deleteError == nil {
				updateExplorerWithCurrentPath()
				explorerView.Update()
			} else {
				logger.Error("fail to delete " + path)
				logger.Error(deleteError.Error())
			}
		})
		explorerView.On("onSelect", func(path string) {
			input.OnSelect(path)
		})
		explorerView.Start()
		<-params.Cancel
	}
}
