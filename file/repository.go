package file

import (
	"io"
)

type Repository interface {
	GetAllNameFileList(dir string) (filesName []string, err error)
	GetFile(dir, name string) (file File, err error)
	CreateFile(dir, name string) (file File, err error)
}

type File interface {
	io.Reader
	io.Writer
	Close() error
	WriteString(s string) (n int, err error)
}
