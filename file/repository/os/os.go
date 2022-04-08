package os

import (
	"fmt"
	"os"

	"github.com/avepa/gethash/file"
)

func (ctrl *Controller) GetAllNameFileList(dir string) (filesName []string, err error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, f := range files {
		if !f.IsDir() {
			filesName = append(filesName, f.Name())
		}
	}
	if len(filesName) == 0 {
		return
	}
	return
}

func (ctrl *Controller) GetFile(dir, name string) (file file.File, err error) {
	path := fmt.Sprintf("%v/%v", dir, name)
	return os.Open(path)
}

func (ctrl *Controller) CreateFile(dir, name string) (file file.File, err error) {
	path := fmt.Sprintf("%v/%v", dir, name)
	return os.Create(path)
}
