package apps

import (
	"github.com/fatih/color"
	react_fullstack_go_server "github.com/shmuelhizmi/react-fullstack-go-server"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/components/desktop"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/types"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/utils"
	"os"
	"path/filepath"
)

func GetNotepadAppInfo() types.RegisteredApp {
	return types.RegisteredApp{
		Name: "Notepad",
		Icon: types.Icon{
			Icon:     "FcFile",
			IconType: "icon",
		},
		NativeIcon: types.NativeIcon{
			Icon:     "filetext1",
			IconType: "AntDesign",
		},
		RegisteredName: "notepad",
		Description:    "just a text editor",
	}
}

func CreateNotepadApp() types.AppRegistrationData {
	info := GetNotepadAppInfo()
	return types.AppRegistrationData{
		Icon:       info.Icon,
		NativeIcon: info.NativeIcon,
		Name:       info.Name,
		DefaultInput: types.NotepadInput{
			FilePath: "",
		},
		Description: info.Description,
		App: func(desktopManager types.DesktopManager, appId int64, input interface{}) (component react_fullstack_go_server.Component) {
			var notepadInput types.NotepadInput = input.(types.NotepadInput)
			window := desktop.CreateWindow(types.CreateWindowParameters{
				Name:  info.Name,
				Title: info.Name,
				Icon:  info.Icon,
				State: types.WindowState{
					Width:     650,
					MaxWidth:  900,
					MinWidth:  550,
					Height:    800,
					MaxHeight: 900,
					MinHeight: 370,
					Position: types.Position{
						X: 100,
						Y: 20,
					},
				},
				DesktopManager: desktopManager,
				App:            NotepadApp(desktopManager, notepadInput),
				AppId:          appId,
			})
			return window.Component
		},
	}
}

func NotepadApp(desktopManager types.DesktopManager, input types.NotepadInput) (component react_fullstack_go_server.Component) {
	return func(params *react_fullstack_go_server.ComponentParams) {
		logger := desktopManager.MountLogger("notepad", color.BgBlue, color.FgBlack)
		logger.Info("running notepad")
		notepadView := params.View(0, types.NotepadViewName, nil)
		currentFilePath := input.FilePath
		currentFileValue := ""
		isSelectingFile := true
		updateNotepadView := func() {
			isSelectingFile = currentFilePath == ""
			notepadView.Params["isSelectingFile"] = isSelectingFile
			if currentFilePath != "" {
				notepadView.Params["name"] = filepath.Base(currentFilePath)
			} else {
				notepadView.Params["name"] = ""
			}
			notepadView.Params["defaultValue"] = currentFileValue
			notepadView.Params["path"] = currentFilePath
		}
		updateNotepadView()
		notepadView.On("onReselectFile", func() {
			currentFilePath = ""
			currentFileValue = ""
			updateNotepadView()
			notepadView.Update()
		})
		notepadView.On("onSave", func(newValue string) {
			currentFileValue = newValue
			updateNotepadView()
			notepadView.Update()
			writeFileError := utils.FSWriteFile(currentFilePath, newValue)
			if writeFileError != nil {
				logger.Error("fail to write file " + currentFilePath)
				logger.Error(writeFileError.Error())
			}
		})
		notepadView.Start()
		getFolderToExplore := func() string {
			if currentFilePath != "" {
				return filepath.Dir(currentFilePath)
			}
			home, getHomeDirError := os.UserHomeDir()
			if getHomeDirError == nil {
				return home
			} else {
				logger.Error("fail to get homedir returning unix default one")
				return "/home/"
			}
		}
		updateFile := func(path string) {
			currentFilePath = path
			fileValue, readFileError := utils.FSReadFile(path)
			if readFileError != nil {
				logger.Error("fail to read file " + path)
				logger.Error(readFileError.Error())
			} else {
				currentFileValue = fileValue
				updateNotepadView()
				notepadView.Update()
			}
		}
		fileSelectComponent := ExplorerApp(desktopManager, types.ExplorerInput{
			Path:         getFolderToExplore(),
			IsCurrentApp: false,
			Type:         types.SelectFile,
			OnSelect:     updateFile,
		})
		params.Run(fileSelectComponent, notepadView)
		<-params.Cancel
	}
}
