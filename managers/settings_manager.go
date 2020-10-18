package managers

import (
	"encoding/json"
	"errors"
	"github.com/fatih/color"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/utils"
	"io/ioutil"
	"os"
	"path"
)

type SettingsManger struct {
	Initialize          func()
	IsInitialized       *bool
	Settings            func() *SettingsObject
	SetSettings         func(newSettings SettingsObject)
	ListenToNewSettings func(listener func(newSettings *SettingsObject))
}

type SettingsManagerDependencies struct {
	logger utils.Logger
}

const (
	SettingsFolder = ".web-desktop-environment-config"
	SettingsFile   = "settings.json"
)

type SettingsObject struct {
	Desktop SettingsObjectDesktop `json:"desktop"`
	Network SettingsObjectNetwork `json:"network"`
}

type SettingsObjectDesktop struct {
	Theme            string `json:"theme"`
	Background       string `json:"background"`
	NativeBackground string `json:"nativeBackground"`
}

type SettingsObjectNetwork struct {
	Ports SettingsObjectNetworkPorts `json:"ports"`
}

type SettingsObjectNetworkPorts struct {
	MainPort  int32 `json:"mainPort"`
	StartPort int32 `json:"startPort"`
	EndPort   int32 `json:"endPort"`
}

func CreateSettingsUpdater(write <-chan SettingsObject, settingsFileFullPath string) {
	for currentWrite := range write {
		fileHandler, fileHandlerCreateError := os.Create(settingsFileFullPath)
		utils.Check(fileHandlerCreateError)
		defaultSettingsJson, _ := json.Marshal(currentWrite)
		_, writeDefaultSettingsError := fileHandler.Write(defaultSettingsJson)
		utils.Check(writeDefaultSettingsError)
		_ = fileHandler.Close()
	}
}

func CreateSettingsManager(dependencies SettingsManagerDependencies) SettingsManger {
	logger := dependencies.logger.Mount("settings manager", color.FgHiMagenta)
	defaultSettings := SettingsObject{
		Desktop: SettingsObjectDesktop{
			Theme:            "transparentDark",
			Background:       "url(https://picsum.photos/id/1039/1920/1080)",
			NativeBackground: "https://picsum.photos/id/237/1080/1920",
		},
		Network: SettingsObjectNetwork{
			Ports: SettingsObjectNetworkPorts{
				MainPort:  5000,
				StartPort: 9200,
				EndPort:   9400,
			},
		},
	}
	settings := defaultSettings
	isInitialized := false
	var newSettingsListeners []func(newSettings *SettingsObject)
	homedir, _ := os.UserHomeDir()
	settingsFolderFullPath := path.Join(homedir, SettingsFolder)
	settingsFileFullPath := path.Join(settingsFolderFullPath, SettingsFile)
	settingsWriteChannel := make(chan SettingsObject)
	go CreateSettingsUpdater(settingsWriteChannel, settingsFileFullPath)
	return SettingsManger{
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
		Settings: func() *SettingsObject {
			if !isInitialized {
				utils.Check(errors.New("trying to read settings before settings manager is initialized"))
			}
			return &settings
		},
		SetSettings: func(newSettings SettingsObject) {
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
		ListenToNewSettings: func(listener func(newSettings *SettingsObject)) {
			newSettingsListeners = append(newSettingsListeners, listener)
		},
	}
}
