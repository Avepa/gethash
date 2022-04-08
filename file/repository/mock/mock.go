package mock

import (
	"fmt"

	"github.com/avepa/gethash/file"
)

var (
	InvelidDirectory = fmt.Errorf("invalid directory")
	InvelidFileName  = fmt.Errorf("invalid file name")
)

type MockRepository struct {
	dir       string
	fileName  string
	filesName []string
	mockFile  map[string]file.File

	errGetFile            error
	errCreatFile          error
	errGetAllNameFileList error
}

func (ctrl *MockRepository) SetDir(dir string) {
	ctrl.dir = dir
	return
}

func (ctrl *MockRepository) SetName(fileName string) {
	ctrl.fileName = fileName
	return
}

func (ctrl *MockRepository) SetMackFile(fileName string, mockFile file.File) {
	ctrl.mockFile = map[string]file.File{fileName: mockFile}
	return
}

func (ctrl *MockRepository) AddMackFile(fileName string, mockFile file.File) {
	if ctrl.mockFile == nil {
		ctrl.mockFile = make(map[string]file.File)
	}
	ctrl.mockFile[fileName] = mockFile
	return
}

func (ctrl *MockRepository) SetFilesName(filesName []string) {
	ctrl.filesName = filesName
	return
}

// func (ctrl *MockRepository) SetFuncGetFile(f func() (file.File, error)) {
// 	ctrl.funcGetFile = f
// 	return
// }

// func (ctrl *MockRepository) SetFuncCreatFile(f func() (file.File, error)) {
// 	ctrl.funcCreatFile = f
// 	return
// }

// func (ctrl *MockRepository) SetFuncDirGetAllNameFileList(f func() ([]string, error)) {
// 	ctrl.funcGetAllNameFileList = f
// 	return
// }

func (ctrl *MockRepository) SetErrGetFile(err error) {
	ctrl.errGetFile = err
	return
}

func (ctrl *MockRepository) SetErrCreatFile(err error) {
	ctrl.errCreatFile = err
	return
}

func (ctrl *MockRepository) SetErrGetAllNameFileList(err error) {
	ctrl.errGetAllNameFileList = err
	return
}

func (ctrl *MockRepository) GetAllNameFileList(dir string) (filesName []string, err error) {
	if ctrl.dir != dir {
		return nil, InvelidDirectory
	}
	if ctrl.errGetAllNameFileList != nil {
		return nil, ctrl.errGetAllNameFileList
	}
	return ctrl.filesName, nil
}

func (ctrl *MockRepository) GetFile(dir, name string) (file file.File, err error) {
	if ctrl.dir != dir {
		return nil, InvelidDirectory
	}
	if _, ok := ctrl.mockFile[name]; !ok {
		return nil, InvelidFileName
	}
	if ctrl.errGetFile != nil {
		return nil, ctrl.errGetFile
	}
	return ctrl.mockFile[name], nil
}

func (ctrl *MockRepository) CreateFile(dir, name string) (file file.File, err error) {
	if ctrl.dir != dir {
		return nil, InvelidDirectory
	}
	if _, ok := ctrl.mockFile[name]; !ok {
		return nil, InvelidFileName
	}
	if ctrl.errCreatFile != nil {
		return nil, ctrl.errCreatFile
	}
	return ctrl.mockFile[name], nil
}

type MockFiles struct {
	close       func() error
	read        func(p []byte) (n int, err error)
	write       func(p []byte) (n int, err error)
	writeString func(s string) (n int, err error)
}

func (ctrl *MockFiles) SetWrite(write func(p []byte) (n int, err error)) {
	ctrl.write = write
}

func (ctrl *MockFiles) SetRead(read func(p []byte) (n int, err error)) {
	ctrl.read = read
}

func (ctrl *MockFiles) SetClose(close func() error) {
	ctrl.close = close
}

func (ctrl *MockFiles) SetWriteString(writeString func(s string) (n int, err error)) {
	ctrl.writeString = writeString
}

func (ctrl *MockFiles) Write(p []byte) (n int, err error) {
	return ctrl.write(p)
}

func (ctrl *MockFiles) Read(p []byte) (n int, err error) {
	return ctrl.read(p)
}

func (ctrl *MockFiles) Close() error {
	return ctrl.close()
}

func (ctrl *MockFiles) WriteString(s string) (n int, err error) {
	return ctrl.writeString(s)
}
