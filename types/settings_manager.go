package types

type SettingsManager struct {
	Initialize          func()
	IsInitialized       *bool
	Settings            func() *SettingsObject
	SetSettings         func(newSettings SettingsObject)
	ListenToNewSettings func(listener func(newSettings *SettingsObject))
}

type SettingsManagerDependencies struct {
	Logger Logger
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
	CustomTheme      Theme  `json:"customTheme"`
}

type SettingsObjectNetwork struct {
	Ports SettingsObjectNetworkPorts `json:"ports"`
}

type SettingsObjectNetworkPorts struct {
	MainPort  int32 `json:"mainPort"`
	StartPort int32 `json:"startPort"`
	EndPort   int32 `json:"endPort"`
}
