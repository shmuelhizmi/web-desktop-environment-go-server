package managers

import (
	"errors"
	"github.com/fatih/color"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/types"
	"net"
	"strconv"
)

func CreatePortManager(dependencies types.PortMangerDependencies) types.PortManager {
	logger := dependencies.Logger.Mount("port manager", color.BgBlue)
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
	return types.PortManager{
		GetDesktopPort: func() (error, int32) {
			mainPort := dependencies.SettingsManager.Settings().Network.Ports.MainPort
			return getAvailablePort(mainPort, mainPort+50)
		},
		GetAppPort: func() (error, int32) {
			portSettings := dependencies.SettingsManager.Settings().Network.Ports
			return getAvailablePort(portSettings.StartPort, portSettings.EndPort)
		},
	}
}
