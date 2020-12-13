package __tests__

import (
	"encoding/json"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/managers"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

func TestSettingsManagerCanGetInitialized(t *testing.T)  {
	desktopManager := managers.CreateDesktopManager()
	desktopManager.SettingsManager.Initialize()
	if *desktopManager.SettingsManager.IsInitialized != true {
		t.Error("Expected settings manager to get initialized")
	}
	t.Log("settings successfully initialized")
}

func ReadFile(path string, to interface{})  {
	fileData, _ := ioutil.ReadFile(path)
	_ = json.Unmarshal(fileData, to)
}

func TestSettingsCanGetUpdated(t *testing.T) {
	desktopManager := managers.CreateDesktopManager()
	desktopManager.SettingsManager.Initialize()
	homedir, _ :=os.UserHomeDir()
	var settings managers.SettingsObject
	ReadFile(path.Join(homedir, managers.SettingsFolder, managers.SettingsFile), &settings)
	theme := settings.Desktop.Theme
	if theme != desktopManager.SettingsManager.Settings().Desktop.Theme {
		t.Error("seems like settings manager was unable to read the real settings")
	}
	desktopManager.SettingsManager.SetSettings(managers.SettingsObject{
		Desktop: managers.SettingsObjectDesktop{
			Theme:            "temp_theme",
			Background:       "__",
			NativeBackground: "__",
		},
		Network: managers.SettingsObjectNetwork{
			Ports: managers.SettingsObjectNetworkPorts{
				MainPort:  0,
				StartPort: 0,
				EndPort:   0,
			},
		},
	})
	time.Sleep(10 * time.Millisecond) // 10 milliseconds should be plenty of time to save the settings to a real file
	var newSettings managers.SettingsObject
	ReadFile(path.Join(homedir, managers.SettingsFolder, managers.SettingsFile), &newSettings)
	if theme == newSettings.Desktop.Theme {
		t.Error("seems like settings manager was unable to save the new settings")
	}
}

func TestSettingsManagerCanSubscribeYouToASaveEvent(t *testing.T)  {
	desktopManager := managers.CreateDesktopManager()
	desktopManager.SettingsManager.Initialize()
	successfullyGotSubscribed := 0
	desktopManager.SettingsManager.ListenToNewSettings(func(newSettings *managers.SettingsObject) {
		successfullyGotSubscribed++
	})
	desktopManager.SettingsManager.ListenToNewSettings(func(newSettings *managers.SettingsObject) {
		successfullyGotSubscribed++
	})
	desktopManager.SettingsManager.ListenToNewSettings(func(newSettings *managers.SettingsObject) {
		successfullyGotSubscribed++
	})
	desktopManager.SettingsManager.SetSettings(managers.SettingsObject{
		Desktop: managers.SettingsObjectDesktop{
			Theme:            "temp_theme",
			Background:       "__",
			NativeBackground: "__",
		},
		Network: managers.SettingsObjectNetwork{
			Ports: managers.SettingsObjectNetworkPorts{
				MainPort:  0,
				StartPort: 0,
				EndPort:   0,
			},
		},
	})
	if successfullyGotSubscribed != 3 {
		t.Error("settings manager did not run all new settings listeners")
	}
}