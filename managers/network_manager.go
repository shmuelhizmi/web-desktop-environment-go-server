package managers

import (
	"github.com/shmuelhizmi/web-desktop-environment-go-server/types"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/utils"
	"net/http"
)

func CreateNetworkManager(dependencies types.NetworkManagerDependencies) types.NetworkManager {
	logger := dependencies.Logger.Mount("network manager")
	server := http.NewServeMux()
	var appIndex uint64 = 0
	var serviceIndex uint64 = 0
	return types.NetworkManager{
		Server: server,
		Initialize: func() {
			port := dependencies.SettingsManager.Settings().Network.Ports.MainPort
			go func() {
				serverStartError := http.ListenAndServe(":"+utils.Int32ToString(port), server)
				if serverStartError != nil {
					logger.Error("fail to start server on Port - " + utils.Int32ToString(port))
					logger.Error(serverStartError.Error())
					panic(serverStartError)
				}
			}()
		},
		GetAppPath: func(appName string) string {
			appIndex++
			return "/apps/" + appName + "_" + utils.UInt64ToString(appIndex) + "/"
		},
		GetServicePath: func(service string) string {
			serviceIndex++
			return "/service/" + service + "_" + utils.UInt64ToString(serviceIndex) + "/"
		},
	}
}
