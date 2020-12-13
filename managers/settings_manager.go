package managers

import (
	"encoding/json"
	"errors"
	"github.com/fatih/color"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/types"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/utils"
	"io/ioutil"
	"os"
	"path"
)

func CreateSettingsUpdater(write <-chan types.SettingsObject, settingsFileFullPath string) {
	for currentWrite := range write {
		fileHandler, fileHandlerCreateError := os.Create(settingsFileFullPath)
		utils.Check(fileHandlerCreateError)
		defaultSettingsJson, _ := json.Marshal(currentWrite)
		_, writeDefaultSettingsError := fileHandler.Write(defaultSettingsJson)
		utils.Check(writeDefaultSettingsError)
		_ = fileHandler.Close()
	}
}

func CreateSettingsManager(dependencies types.SettingsManagerDependencies) types.SettingsManager {
	logger := dependencies.Logger.Mount("settings manager", color.FgHiMagenta)
	defaultSettings := types.SettingsObject{
		Desktop: types.SettingsObjectDesktop{
			Theme:            "transparentDark",
			Background:       "url(https://picsum.photos/id/1039/1920/1080)",
			NativeBackground: "https://picsum.photos/id/237/1080/1920",
		},
		Network: types.SettingsObjectNetwork{
			Ports: types.SettingsObjectNetworkPorts{
				MainPort:  5000,
			},
		},
	}
	settings := defaultSettings
	isInitialized := false
	var newSettingsListeners []func(newSettings *types.SettingsObject)
	homedir, _ := os.UserHomeDir()
	settingsFolderFullPath := path.Join(homedir, types.SettingsFolder)
	settingsFileFullPath := path.Join(settingsFolderFullPath, types.SettingsFile)
	settingsWriteChannel := make(chan types.SettingsObject)
	go CreateSettingsUpdater(settingsWriteChannel, settingsFileFullPath)
	return types.SettingsManager{
		Initialize: func() {
			_, folderStatErr := os.Stat(settingsFolderFullPath)
			utils.Check(folderStatErr)
			if os.IsNotExist(folderStatErr) {
				makeSettingsDirError := os.Mkdir(settingsFolderFullPath, 0700)
				utils.Check(makeSettingsDirError)
			}
			_, fileStatErr := os.Stat(settingsFileFullPath)
			if os.IsNotExist(fileStatErr) {
				settingsWriteChannel <- settings
			} else {
				fileData, settingsReadError := ioutil.ReadFile(settingsFileFullPath)
				utils.Check(settingsReadError)
				settingsJsonParseError := json.Unmarshal(fileData, &settings)
				if settingsJsonParseError != nil {
					logger.Error("fail to load settings file - creating default one \n error is: " + settingsJsonParseError.Error())
					settingsWriteChannel <- settings
				}
			}
			logger.Info("finish initializing settings manager")
			isInitialized = true
		},
		IsInitialized: &isInitialized,
		Settings: func() *types.SettingsObject {
			if !isInitialized {
				utils.Check(errors.New("trying to read settings before settings manager is initialized"))
			}
			return &settings
		},
		SetSettings: func(newSettings types.SettingsObject) {
			if !isInitialized {
				utils.Check(errors.New("trying to read settings before settings manager is initialized"))
			}
			settings = newSettings
			go func() {
				logger.Info("set settings start")
				settingsWriteChannel <- newSettings
				logger.Info("set settings end")
			}()
			for _, listener := range newSettingsListeners {
				listener(&settings)
			}
		},
		ListenToNewSettings: func(listener func(newSettings *types.SettingsObject)) {
			newSettingsListeners = append(newSettingsListeners, listener)
		},
	}
}
