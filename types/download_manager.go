package types

type DownloadManagerDependencies struct {
	Logger      Logger
	PortManager PortManager
}

type DownloadManager struct {
	AddFile    func(path string) (hash string)
	Port       *int32
	Initialize func() (initializationError error)
}
