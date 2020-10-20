package utils

import (
	"github.com/shmuelhizmi/web-desktop-environment-go-server/types"
	"io/ioutil"
)

func ListFilesInDir(dir string) []types.ExplorerFile {
	var files []types.ExplorerFile
	readResult, _ := ioutil.ReadDir(dir)
	for _, file := range readResult {
		files = append(files, types.ExplorerFile{
			IsFolder: file.IsDir(),
			Name:     file.Name(),
		})
	}

	return files
}
