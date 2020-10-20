package types

type ExplorerViewType string

const (
	Explore      ExplorerViewType = "explorer"
	SelectFile   ExplorerViewType = "select-file"
	SelectFolder ExplorerViewType = "select-folder"
)

type ExplorerInput struct {
	Path         string
	IsCurrentApp bool
	Type         ExplorerViewType
	onSelect     func(path string)
}

type ExplorerFile struct {
	IsFolder bool   `json:"isFolder"`
	Name     string `json:"name"`
}
