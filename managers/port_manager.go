package managers

import (
	"errors"
	"github.com/fatih/color"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/utils"
	"net"
	"strconv"
)

type PortMangerDependencies struct {
	logger utils.Logger
	settingsManager SettingsManger
}

type PortManager struct {
	getDesktopPort func() (error, int32)
	getAppPort     func() (error, int32)
}

func CreatePortManager(dependencies PortMangerDependencies) PortManager {
	logger := dependencies.logger.Mount("port manager", color.BgBlue)
	isPortAvailable := func(port string) bool {
		ln, err := net.Listen("tcp", ":"+port)
		if err != nil {
			return false
		}
		_ = ln.Close()
		return true
	}
	getAvailablePort := func(port int32, maxPort int32) (error, int32) {
		currentPort := port
		for currentPort <= maxPort {
			stringPort := strconv.FormatInt(int64(currentPort), 10)
			if isPortAvailable(stringPort) {
				logger.Info("found port - " + stringPort)
				return nil, currentPort
			}
			currentPort++
		}
		return errors.New("no port available in range"), -1
	}
	return PortManager{
		getDesktopPort: func() (error, int32) {
			mainPort := dependencies.settingsManager.Settings().Network.Ports.MainPort
			return getAvailablePort(mainPort, mainPort + 50)
		},
		getAppPort: func() (error, int32) {
			portSettings := dependencies.settingsManager.Settings().Network.Ports
			return getAvailablePort(portSettings.StartPort, portSettings.EndPort)
		},
	}
}
