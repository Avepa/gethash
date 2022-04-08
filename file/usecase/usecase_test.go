package usecase

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/avepa/gethash/file"
	"github.com/avepa/gethash/file/repository/mock"
	tHash "github.com/avepa/gethash/tools/hash"
	"golang.org/x/exp/slices"
)

// вернуть ошибку на получении списка файлов
// вернуть ошибку на сохранении файлов
func TestStartCountingHashes1(t *testing.T) {
	ctrlRepo := &mock.MockRepository{}
	ctrl := NewController(ctrlRepo)
	cfg := &file.Config{
		FolderDir:    "main/files",
		SaveFileDir:  "main/files",
		SaveFileName: "hash.txt",
		HashName:     tHash.Keccak256.String(),
		MaxProc:      2,
	}
	ctrlRepo.SetDir("main/files")

	type typeTestData struct {
		fileName string
		fileData []byte
	}

	close := func() (err error) { return }
	filesName := []string{"zero.txt", "one.txt", "two.txt", "three.txt", "four.txt", "five.txt", "six.txt", "seven.txt", "eight.txt", "nine.txt"}
	files := []typeTestData{
		{
			fileName: "eight.txt",
			fileData: []byte("eight eight eight eight eight eight"),
		},
		{
			fileName: "five.txt",
			fileData: []byte("five five five five five five five"),
		},
		{
			fileName: "four.txt",
			fileData: []byte("four four four four four four four"),
		},
		{
			fileName: "nine.txt",
			fileData: []byte("nine nine nine nine nine nine nine"),
		},
		{
			fileName: "one.txt",
			fileData: []byte("one one one one one one one one"),
		},
		{
			fileName: "seven.txt",
			fileData: []byte("seven seven seven seven seven seven"),
		},
		{
			fileName: "six.txt",
			fileData: []byte("sixe sixe sixe sixe sixe sixe sixe"),
		},
		{
			fileName: "three.txt",
			fileData: []byte("three three three three three three"),
		},
		{
			fileName: "two.txt",
			fileData: []byte("two two two two two two two two"),
		},
		{
			fileName: "zero.txt",
			fileData: []byte("zero zero zero zero zero zero zero zero"),
		},
	}

	for _, f := range files {
		mf := &mock.MockFiles{}
		mf.SetClose(close)
		mf.SetRead(bytes.NewReader(f.fileData).Read)
		ctrlRepo.AddMackFile(f.fileName, mf)
	}

	ctrlRepo.SetFilesName(filesName)

	ws := func(s string) (n int, err error) {
		if len(files) > 0 {
			h, err := tHash.StringToHashes(cfg.HashName)
			if err != nil {
				t.Errorf("This code does not exist, hash %v", cfg.HashName)
			}
			newHash, err := h.NewHash()
			if err != nil {
				t.Error("failed to get hash function")
			}

			hasher := newHash()
			_, err = hasher.Write(files[0].fileData)
			if err != nil {
				t.Error("failed to get hash function")
			}

			data := fmt.Sprintf("%v\n", hex.EncodeToString(hasher.Sum(nil)))
			if s != data {
				t.Errorf("Invalid hash received. File name: %v, expected: %v, received: %v.",
					files[0].fileName, data, s)
			}
			files = files[1:]
		} else {
			err = errors.New("hash limit exceeded")
			t.Error(err)
		}
		return
	}

	mfWrite := &mock.MockFiles{}
	mfWrite.SetClose(close)
	mfWrite.SetWriteString(ws)
	ctrlRepo.AddMackFile(cfg.SaveFileName, mfWrite)

	err := ctrl.StartCountingHashes(context.Background(), cfg)
	if err != nil {
		t.Errorf("Unexpected error, err: %v.", err)
	}
	if len(files) > 0 {
		t.Error("not all files processed")
	}

	errInput := errors.New("error creatfile")
	ctrlRepo.SetErrCreatFile(errInput)
	err = ctrl.StartCountingHashes(context.Background(), cfg)
	if err != errInput {
		t.Errorf("Unexpected error, err: %v.", err)
	}

	errInput = errors.New("error gGetAllNameFileList")
	ctrlRepo.SetErrGetAllNameFileList(errInput)
	err = ctrl.StartCountingHashes(context.Background(), cfg)
	if err != errInput {
		t.Errorf("Unexpected error, err: %v.", err)
	}

	cfg.HashName = "SHA3-2048"
	err = ctrl.StartCountingHashes(context.Background(), cfg)
	if err != tHash.NotFound {
		t.Errorf("Unexpected error, err: %v.", err)
	}
	return
}

func TestCountHash(t *testing.T) {
	ctrlRepo := &mock.MockRepository{}
	mockFile := &mock.MockFiles{}
	mockFile.SetClose(func() error { return nil })
	hashFunc, err := tHash.MD5.NewHash()
	if err != nil {
		t.Fatal(err)
	}

	type inputData struct {
		fileDir  string
		fileName string
		fileData string
	}

	type StructTestData struct {
		input        inputData
		expectedData HashFile
	}

	testData := []StructTestData{
		{
			input: inputData{
				fileName: "zero.txt",
				fileData: "zero, zero",
			},
			expectedData: HashFile{
				FileName: "one.txt",
				Error:    mock.InvelidDirectory,
			},
		},
		{
			input: inputData{
				fileDir:  "main/dir",
				fileName: "one.txt",
				fileData: "one one one",
			},
			expectedData: HashFile{
				FileName: "one.txt",
			},
		},
		{
			input: inputData{
				fileDir:  "main/dir",
				fileName: "two.txt",
				fileData: "two two two",
			},
			expectedData: HashFile{
				FileName: "two.txt",
			},
		},
	}

	ctrl := &ControllerFileHash{
		nameFile: make(chan string),
		hashFile: make(chan HashFile),
		newHash:  hashFunc,
		ctrlRepo: ctrlRepo,
	}

	var ctx context.Context
	ctx, ctrl.chCancel = context.WithCancel(context.Background())

	ctrl.wgCountHash.Add(1)
	go ctrl.countHash(ctx, "main/dir")

	for _, td := range testData {
		reader := bytes.NewReader([]byte(td.input.fileData))
		mockFile.SetRead(reader.Read)
		ctrlRepo.SetDir(td.input.fileDir)
		ctrlRepo.SetName(td.input.fileName)
		ctrlRepo.AddMackFile(td.input.fileName, mockFile)

		ctrl.nameFile <- td.input.fileName
		hashData := <-ctrl.hashFile
		if hashData.FileName != td.input.fileName {
			t.Errorf("Names don't match. Expected: %v, received: %v.",
				td.input.fileData, hashData.FileName)
			continue
		}
		if !errors.Is(hashData.Error, td.expectedData.Error) {
			t.Errorf("Got wrong error. Expected: %v, received: %v.",
				td.expectedData.Error, hashData.Error)
			continue
		}
		if hashData.Error == nil {
			hasher := hashFunc()
			hasher.Write([]byte(td.input.fileData))
			if !slices.Equal(hashData.FileHash, hasher.Sum(nil)) {
				t.Errorf("Hashes do not match. Expected: %v, received: %v.",
					hex.EncodeToString(hasher.Sum(nil)), hex.EncodeToString(hashData.FileHash))
			}
		}
	}

	ctrl.StopCountHash()
	return
}

func TestSaveHash(t *testing.T) {
	ctrlRepo := &mock.MockRepository{}
	ctrl := &ControllerFileHash{
		hashFile: make(chan HashFile, 10),
		ctrlRepo: ctrlRepo,
	}

	filesName := []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	hashFile := []HashFile{
		{
			FileName: "zero",
			FileHash: []byte("0"),
		},
		{
			FileName: "one",
			Error:    errors.New("1"),
		},
		{
			FileName: "two",
			FileHash: []byte("2"),
		},
		{
			FileName: "three",
			FileHash: []byte("3"),
		},
		{
			FileName: "four",
			FileHash: []byte("4"),
		},
		{
			FileName: "five",
			Error:    errors.New("5"),
		},
		{
			FileName: "six",
			Error:    errors.New("6"),
		},
		{
			FileName: "seven",
			FileHash: []byte("7"),
		},
		{
			FileName: "eight",
			FileHash: []byte("8"),
		},
		{
			FileName: "nine",
			FileHash: []byte("9"),
		},
	}

	for i := len(hashFile) / 2; i < len(hashFile); i++ {
		ctrl.hashFile <- hashFile[i]
	}
	for i := 0; i < len(hashFile)/2; i++ {
		ctrl.hashFile <- hashFile[i]
	}

	var ctx context.Context
	ctx, ctrl.shCancel = context.WithCancel(context.Background())
	err := ctrl.saveHash(ctx, "main/data", "hash.txt", filesName)
	if err != mock.InvelidDirectory {
		t.Errorf("Got wrong error. Expected: %v, received: %v.",
			mock.InvelidDirectory, err)
	}

	ctrlRepo.SetDir("main/data")
	err = ctrl.saveHash(ctx, "main/data", "hash.txt", filesName)
	if err != mock.InvelidFileName {
		t.Errorf("Got wrong error. Expected: %v, received: %v.",
			mock.InvelidFileName, err)
	}

	errInput := errors.New("error")
	ctrlRepo.AddMackFile("hash.txt", nil)
	ctrlRepo.SetName("hash.txt")
	ctrlRepo.SetErrCreatFile(errInput)
	err = ctrl.saveHash(ctx, "main/data", "hash.txt", filesName)
	if !errors.Is(errInput, err) {
		t.Errorf("Got wrong error. Expected: %v, received: %v.",
			errInput, err)
	}

	ws := func(s string) (n int, err error) {
		if len(hashFile) > 0 {
			if hashFile[0].Error != nil {
				if s != fmt.Sprintf("Error: %v\n", hashFile[0].Error.Error()) {
					t.Errorf("Got wrong error. File name: %v, expected: %v, received: %v.",
						hashFile[0].FileName, hashFile[0].Error, s)
				}
			} else {
				data := fmt.Sprintf("%v\n", hex.EncodeToString(hashFile[0].FileHash))
				if s != data {
					t.Errorf("Invalid hash received. File name: %v, expected: %v, received: %v.",
						hashFile[0].FileName, s, data)
				}
			}
			hashFile = hashFile[1:]
		} else {
			err = errors.New("hash limit exceeded")
			t.Error(err)
		}
		return
	}

	mockFile := &mock.MockFiles{}
	mockFile.SetClose(func() error { return nil })
	mockFile.SetWriteString(ws)
	ctrlRepo.SetMackFile("hash.txt", mockFile)

	ctrl.wgSaveHash.Add(1)
	ctrlRepo.SetErrCreatFile(nil)
	go ctrl.saveHash(ctx, "main/data", "hash.txt", filesName)
	ctrl.wgSaveHash.Wait()
	if len(hashFile) > 0 {
		t.Fatal("Received more response than expected.")
	}

	ctrl.wgSaveHash.Add(1)
	go ctrl.saveHash(ctx, "main/data", "hash.txt", filesName)
	ctrl.shCancel()
	return
}

func TestStopCountHash(t *testing.T) {
	ctrl := &ControllerFileHash{
		nameFile: make(chan string),
	}
	ctrl.wgCountHash.Add(1)

	var ctx context.Context
	ctx, ctrl.chCancel = context.WithCancel(context.Background())
	go func() {
		<-ctx.Done()
		time.Sleep(time.Second / 2)
		ctrl.wgCountHash.Done()
	}()

	ctrl.StopCountHash()
	if _, open := <-ctrl.nameFile; open {
		t.Fatal("nameFile channel must be closed")
	}
	return
}

func TestStopSaveHash(t *testing.T) {
	ctrl := &ControllerFileHash{
		hashFile: make(chan HashFile),
	}
	ctrl.wgSaveHash.Add(1)

	var ctx context.Context
	ctx, ctrl.shCancel = context.WithCancel(context.Background())
	go func() {
		<-ctx.Done()
		time.Sleep(time.Second / 2)
		ctrl.wgSaveHash.Done()
	}()

	ctrl.StopSaveHash()
	if _, open := <-ctrl.hashFile; open {
		t.Fatal("hashFile channel must be closed")
	}
	return
}
