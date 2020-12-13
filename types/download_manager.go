package types

type DownloadManagerDependencies struct {
	Logger         Logger
	NetworkManager NetworkManager
}

type DownloadManager struct {
	AddFile    func(path string) (hash string)
	Path       *string
	Initialize func() (initializationError error)
}
