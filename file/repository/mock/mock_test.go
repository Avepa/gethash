package mock

import (
	"bytes"
	"errors"
	"testing"

	"golang.org/x/exp/slices"
)

func TestSetDir(t *testing.T) {
	mr := &MockRepository{}
	mr.SetDir("main")
	if mr.dir != "main" {
		t.Error("Failed to set directory.")
	}
	return
}

func TestSetName(t *testing.T) {
	mr := &MockRepository{}
	mr.SetName("hash.txt")
	if mr.fileName != "hash.txt" {
		t.Error("Failed to set name.")
	}
	return
}

func TestSetMackFile(t *testing.T) {
	mr := &MockRepository{}
	mf := &MockFiles{}
	mr.SetMackFile("test.txt", mf)
	if mr.mockFile == nil || len(mr.mockFile) != 1 || mr.mockFile["test.txt"] != mf {
		t.Error("Failed to set mockFile")
	}
	return
}

func TestAddMackFile(t *testing.T) {
	mr := &MockRepository{}
	mf := &MockFiles{}

	mr.AddMackFile("test.txt", mf)
	mr.AddMackFile("test1.txt", mf)
	if mr.mockFile == nil || len(mr.mockFile) != 2 || mr.mockFile["test.txt"] != mf {
		t.Error("Failed to set mockFile")
	}
	return
}

func TestSetFilesName(t *testing.T) {
	mr := &MockRepository{}
	fn := []string{"1.txt", "2.txt", "3.txt", "4.txt", "5.txt"}
	mr.SetFilesName(fn)
	if !slices.Equal(mr.filesName, fn) {
		t.Error("Failed to set filesName")
	}
	return
}

func TestSetErrGetFile(t *testing.T) {
	mr := &MockRepository{}
	err := errors.New("error")
	mr.SetErrGetFile(err)
	if err != mr.errGetFile {
		t.Error("Failed to set errGetFile")
	}
	return
}

func TestSetErrCreatFile(t *testing.T) {
	mr := &MockRepository{}
	err := errors.New("error")
	mr.SetErrCreatFile(err)
	if err != mr.errCreatFile {
		t.Error("Failed to set errCreatFile")
	}
	return
}

func TestSetErrGetAllNameFileList(t *testing.T) {
	mr := &MockRepository{}
	err := errors.New("error")
	mr.SetErrGetAllNameFileList(err)
	if err != mr.errGetAllNameFileList {
		t.Error("Failed to set errGetAllNameFileList")
	}
	return
}

func TestGetAllNameFileList(t *testing.T) {
	mr := &MockRepository{}
	fn, err := mr.GetAllNameFileList("main")
	if err != InvelidDirectory {
		t.Errorf("Received an unexpected error. Expected: %v, received: %v.",
			InvelidDirectory, err)
	}
	if fn != nil {
		t.Error("function must return nil")
	}

	filesName := []string{"1.txt", "2.txt", "3.txt", "4.txt", "5.txt"}
	mr.SetFilesName(filesName)

	fn, err = mr.GetAllNameFileList("")
	if err != nil {
		t.Errorf("Received an unexpected error. Expected: %v, received: %v.",
			nil, err)
	}
	if !slices.Equal(fn, filesName) {
		t.Errorf("Arrays are not equal. Expected: %v, received: %v.", filesName, fn)
	}

	errInput := errors.New("error")
	mr.SetErrGetAllNameFileList(errInput)
	fn, err = mr.GetAllNameFileList("main")
	if err == errInput {
		t.Errorf("Received an unexpected error. Expected: %v, received: %v.",
			errInput, err)
	}
	if fn != nil {
		t.Error("function must return nil")
	}
	return
}

func TestGetFile(t *testing.T) {
	mr := &MockRepository{}
	fn, err := mr.GetFile("main", "main")
	if err != InvelidDirectory {
		t.Errorf("Received an unexpected error. Expected: %v, received: %v.",
			InvelidDirectory, err)
	}
	if fn != nil {
		t.Error("function must return nil")
	}

	fn, err = mr.GetFile("", "main")
	if err != InvelidFileName {
		t.Errorf("Received an unexpected error. Expected: %v, received: %v.",
			InvelidFileName, err)
	}
	if fn != nil {
		t.Error("function must return nil")
	}

	mf := &MockFiles{}
	mr.SetMackFile("", mf)

	fn, err = mr.GetFile("", "")
	if err != nil {
		t.Errorf("Received an unexpected error. Expected: %v, received: %v.",
			nil, err)
	}
	if _, ok := fn.(*MockFiles); !ok {
		t.Errorf("Returned variable with wrong type. Expected: %T, received: %T.", mf, fn)
	}

	errInput := errors.New("error")
	mr.SetErrGetFile(errInput)
	fn, err = mr.GetFile("", "")
	if err != errInput {
		t.Errorf("Received an unexpected error. Expected: %v, received: %v.",
			errInput, err)
	}
	if fn != nil {
		t.Error("function must return nil")
	}
	return
}

func TestCreateFile(t *testing.T) {
	mr := &MockRepository{}
	fn, err := mr.CreateFile("main", "main")
	if err != InvelidDirectory {
		t.Errorf("Received an unexpected error. Expected: %v, received: %v.",
			InvelidDirectory, err)
	}
	if fn != nil {
		t.Error("function must return nil")
	}

	fn, err = mr.CreateFile("", "main")
	if err != InvelidFileName {
		t.Errorf("Received an unexpected error. Expected: %v, received: %v.",
			InvelidFileName, err)
	}
	if fn != nil {
		t.Error("function must return nil")
	}

	mf := &MockFiles{}
	mr.SetMackFile("", mf)

	fn, err = mr.CreateFile("", "")
	if err != nil {
		t.Errorf("Received an unexpected error. Expected: %v, received: %v.",
			nil, err)
	}
	if _, ok := fn.(*MockFiles); !ok {
		t.Errorf("Returned variable with wrong type. Expected: %T, received: %T.", mf, fn)
	}

	errInput := errors.New("error")
	mr.SetErrGetFile(errInput)
	fn, err = mr.CreateFile("", "")
	if err == errInput {
		t.Errorf("Received an unexpected error. Expected: %v, received: %v.",
			errInput, err)
	}
	if fn == nil {
		t.Error("function must return nil")
	}
	return
}

func TestSetWrite(t *testing.T) {
	mf := &MockFiles{}
	mf.SetWrite(bytes.NewBuffer([]byte{}).Write)
	if mf.write == nil {
		t.Error("Failed to set write.")
	}
	return
}

func TestSetRead(t *testing.T) {
	mf := &MockFiles{}
	mf.SetRead(bytes.NewBuffer([]byte{}).Read)
	if mf.read == nil {
		t.Error("Failed to set read.")
	}
	return
}

func TestSetClose(t *testing.T) {
	mf := &MockFiles{}
	close := func() (err error) { return }
	mf.SetClose(close)
	if mf.close == nil {
		t.Error("Failed to set close.")
	}
	return
}

func TestSetWriteString(t *testing.T) {
	mf := &MockFiles{}
	mf.SetWriteString(bytes.NewBuffer([]byte{}).WriteString)
	if mf.writeString == nil {
		t.Error("Failed to set writeString.")
	}
	return
}

func TestWrite(t *testing.T) {
	mf := &MockFiles{}
	i := 0
	writer := func(p []byte) (n int, err error) {
		i++
		return
	}

	mf.SetWrite(writer)
	_, _ = mf.Write(nil)
	if i > 1 {
		t.Error("Called more than once")
	} else if i < 1 {
		t.Error("The function was not called at all")
	}
	return
}

func TestRead(t *testing.T) {
	mf := &MockFiles{}
	i := 0
	Read := func(p []byte) (n int, err error) {
		i++
		return
	}

	mf.SetRead(Read)
	_, _ = mf.Read(nil)
	if i > 1 {
		t.Error("Called more than once")
	} else if i < 1 {
		t.Error("The function was not called at all")
	}
	return
}

func TestWriteString(t *testing.T) {
	mf := &MockFiles{}
	i := 0
	writeString := func(p string) (n int, err error) {
		i++
		return
	}

	mf.SetWriteString(writeString)
	_, _ = mf.WriteString("")
	if i > 1 {
		t.Error("Called more than once")
	} else if i < 1 {
		t.Error("The function was not called at all")
	}
	return
}

func TestClose(t *testing.T) {
	mf := &MockFiles{}
	i := 0
	close := func() (err error) {
		i++
		return
	}

	mf.SetClose(close)
	mf.Close()
	if i > 1 {
		t.Error("Called more than once")
	} else if i < 1 {
		t.Error("The function was not called at all")
	}
	return
}
