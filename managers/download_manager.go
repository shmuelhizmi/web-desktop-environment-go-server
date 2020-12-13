package managers

import (
	"github.com/fatih/color"
	"github.com/gorilla/mux"
	react_fullstack_go_server "github.com/shmuelhizmi/react-fullstack-go-server"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/types"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/utils"
	"net/http"
	path2 "path"
)

func CreateDownloadManager(dependencies types.DownloadManagerDependencies) (downloadManager types.DownloadManager) {
	logger := dependencies.Logger.Mount("download manager", color.FgRed)
	hashToFiles := make(map[string]string, 0)
	var servicePath string
	return types.DownloadManager{
		AddFile: func(path string) string {
			uuid := react_fullstack_go_server.StringUuid()
			hash := uuid + "." + path2.Ext(path)
			hashToFiles[hash] = path
			return hash
		},
		Path: &servicePath,
		Initialize: func() (initializationError error) {
			servicePath = dependencies.NetworkManager.GetServicePath("downloadManager")
			muxServer := dependencies.NetworkManager.Server
			muxServer.HandleFunc(servicePath+"{hash}", func(writer http.ResponseWriter, request *http.Request) {
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
							logger.Warn("user try to download file that does not exist")
						}
					}
				}
			})
			return nil
		},
	}
}
