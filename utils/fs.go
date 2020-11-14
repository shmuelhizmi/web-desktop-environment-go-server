package utils

import (
	"fmt"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/types"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
)

func FSListFilesInDir(dir string) []types.ExplorerFile {
	var files = make([]types.ExplorerFile, 0)
	readResult, _ := ioutil.ReadDir(dir)
	for _, file := range readResult {
		files = append(files, types.ExplorerFile{
			IsFolder: file.IsDir(),
			Name:     file.Name(),
		})
	}

	return files
}

func FSCopyDirectory(scrDir, dest string) error {
	entries, err := ioutil.ReadDir(scrDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(scrDir, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		stat, ok := fileInfo.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("failed to get raw syscall.Stat_t data for '%s'", sourcePath)
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := FSCreateIfNotExists(destPath, 0755); err != nil {
				return err
			}
			if err := FSCopyDirectory(sourcePath, destPath); err != nil {
				return err
			}
		case os.ModeSymlink:
			if err := FSCopySymLink(sourcePath, destPath); err != nil {
				return err
			}
		default:
			if err := FSCopyFile(sourcePath, destPath); err != nil {
				return err
			}
		}

		if err := os.Lchown(destPath, int(stat.Uid), int(stat.Gid)); err != nil {
			return err
		}

		isSymlink := entry.Mode()&os.ModeSymlink != 0
		if !isSymlink {
			if err := os.Chmod(destPath, entry.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

func FSCopyFile(srcFile, dstFile string) error {
	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}

	defer out.Close()

	in, err := os.Open(srcFile)
	defer in.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

func FSExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func FSCreateIfNotExists(dir string, perm os.FileMode) error {
	if FSExists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}

func FSCopySymLink(source, dest string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return err
	}
	return os.Symlink(link, dest)
}

func FSCopy(path string, newPath string) error {
	fileOrFolder, err := os.Stat(path)
	if err != nil {
		return err
	}
	switch fileOrFolder.Mode() {
	case os.ModeDir:
		return FSCopyDirectory(path, newPath)
	case os.ModeSymlink:
		return FSCopySymLink(path, newPath)
	default:
		return FSCopyFile(path, newPath)
	}
}

func FSCreateEmptyFile(path string) error {
	return ioutil.WriteFile(path, nil, 0)
}

func FSWriteFile(path string, value string) error {
	return ioutil.WriteFile(path, []byte(value), 0)
}

func FSReadFile(path string) (string, error) {
	fileBytes, readFileError := ioutil.ReadFile(path)
	if readFileError != nil {
		return "", readFileError
	}
	return string(fileBytes), nil
}

func FSCreateFolder(path string) error {
	_, statError := os.Stat(path)

	if os.IsNotExist(statError) {
		mkdirError := os.MkdirAll(path, 0755)
		return mkdirError

	}
	return statError
}

func FSDelete(path string) error {
	return os.RemoveAll(path)
}
