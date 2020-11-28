package managers

import (
	"errors"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
	react_fullstack_go_server "github.com/shmuelhizmi/react-fullstack-go-server"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/types"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/utils"
	"log"
	"net/http"
	path2 "path"
)

func CreateDownloadManager(dependencies types.DownloadManagerDependencies) (downloadManager types.DownloadManager) {
	logger := dependencies.Logger.Mount("download manager", color.FgRed)
	var managerPort int32 = 0
	hashToFiles := make(map[string]string, 0)
	return types.DownloadManager{
		AddFile: func(path string) string {
			uuid := react_fullstack_go_server.StringUuid()
			hash := uuid + "." + path2.Ext(path)
			hashToFiles[hash] = path
			return hash
		},
		Port: &managerPort,
		Initialize: func() (initializationError error) {
			getPortError, newManagerPort := dependencies.PortManager.GetAppPort()
			if getPortError != nil {
				logger.Error("fail to initialize download manager: fail to get available port")
				logger.Error(getPortError.Error())
				return errors.New("fail to initialize download manager: fail to get available port")
			}
			managerPort = newManagerPort
			go func() {
				muxDownloadServer := mux.NewRouter()
				muxDownloadServer.HandleFunc("/{hash}", func(writer http.ResponseWriter, request *http.Request) {
					params := mux.Vars(request)
					hash, isHashFound := params["hash"]
					if isHashFound {
						path, isPathFound := hashToFiles[hash]
						if isPathFound {
							fileStr, readFileError := utils.FSReadFile(path)
							if readFileError == nil {
								writer.WriteHeader(200)
								writer.Header().Set("Content-Disposition", "attachment")
								writer.Header().Set("filename", path2.Base(path))
								writer.Write([]byte(fileStr))

							} else {
								writer.WriteHeader(500)
								writer.Write([]byte("unable to find file in file system"))
							}
						}
					}
				})
				log.Panic(http.ListenAndServe(":"+utils.Int32ToString(managerPort), muxDownloadServer))
			}()
			return nil
		},
	}
}
