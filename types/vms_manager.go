package types

import react_fullstack_go_server "github.com/shmuelhizmi/react-fullstack-go-server"

type VMSManagerDependencies struct {
	Logger              Logger
	PortManager         PortManager
	SettingsManager     SettingsManager
	ApplicationsManager ApplicationsManager
}

type VMPermission string

const (
	MathPermission   VMPermission = "math"
	OSPermission     VMPermission = "os"
	TextPermission   VMPermission = "text"
	TimePermission   VMPermission = "times"
	RandomPermission VMPermission = "rand"
	FmtPermission    VMPermission = "fmt"
	JSONPermission   VMPermission = "json"
	Base64Permission VMPermission = "base64"
	HexPermission    VMPermission = "hex"
)

type VMSManager struct {
	RunVMApp func(script string, permissions []VMPermission, appParams *react_fullstack_go_server.ComponentParams) error
}
