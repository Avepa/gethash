package usecase

import (
	"context"
	"hash"
	"sync"

	"github.com/avepa/gethash/file"
	tHash "github.com/avepa/gethash/tools/hash"
)

type Controller struct {
	ctrlRepo file.Repository
}

func NewController(ctrl file.Repository) *Controller {
	return &Controller{
		ctrlRepo: ctrl,
	}
}

type HashFile struct {
	FileName string
	FileHash []byte
	Error    error
}

type ControllerFileHash struct {
	nameFile chan string
	hashFile chan HashFile
	newHash  func() hash.Hash

	shCancel    context.CancelFunc
	chCancel    context.CancelFunc
	wgSaveHash  sync.WaitGroup
	wgCountHash sync.WaitGroup
	ctrlRepo    file.Repository
}

func newControllerFileHash(maxProc int, nameHash string, ctrlRepo file.Repository) (*ControllerFileHash, error) {
	ctrl := &ControllerFileHash{
		nameFile: make(chan string, maxProc),
		hashFile: make(chan HashFile, maxProc),
		ctrlRepo: ctrlRepo,
	}

	ht, err := tHash.StringToHashes(nameHash)
	if err != nil {
		return nil, err
	}
	ctrl.newHash, err = ht.NewHash()
	if err != nil {
		return nil, err
	}

	return ctrl, err
}
