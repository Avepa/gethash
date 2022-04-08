package os

import (
	"flag"
	"os"
	"testing"
)

func getTestDir() (dir string, err error) {
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	f.StringVar(&dir, "test_dir", "testdata/filetest", "директория для тестовых данных")
	err = f.Parse(os.Args)
	return
}

func TestGetAllNameFileList(t *testing.T) {
	nFile := 2
	dir, err := getTestDir()
	if err != nil {
		t.Fatalf("Failed to get test directory, err: %v.", err)
	}

	ctrl := NewController()
	filesName, err := ctrl.GetAllNameFileList(dir)
	if err != nil {
		t.Error(err)
	}
	if len(filesName) != nFile {
		t.Fatal(filesName)
	}

	filesName, err = ctrl.GetAllNameFileList(dir)
	if err == nil {
		t.Error(err)
	}
	if len(filesName) > 0 {
		t.Error(err)
	}
}

func TestCreateFile(t *testing.T) {
	dir, err := getTestDir()
	if err != nil {
		t.Fatalf("Failed to get test directory, err: %v.", err)
	}

	ctrl := NewController()
	files, err := os.ReadDir(dir)
	for _, f := range files {
		if !f.IsDir() {
			_, err := ctrl.CreateFile(dir, "")
			if err == nil {
				t.Errorf("")
			}
		}
	}
}

func TestGetFile(t *testing.T) {
	dir, err := getTestDir()
	if err != nil {
		t.Fatalf("Failed to get test directory, err: %v.", err)
	}

	ctrl := NewController()
	files, err := os.ReadDir(dir)
	for _, f := range files {
		if !f.IsDir() {
			file, err := ctrl.GetFile(dir, f.Name())
			if err != nil {
				t.Errorf("")
			}
			if file == nil {
				t.Errorf("")
			}
			file.Close()

		}
	}
}
