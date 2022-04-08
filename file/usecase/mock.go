package usecase

import (
	"context"
	"errors"

	"github.com/avepa/gethash/file"
)

var (
	СonfigMissing = errors.New("Сonfig missing.")
)

type MockUsecase struct {
	err error
}

func (ctrl *MockUsecase) SetError(err error) {
	ctrl.err = err
	return
}

func (ctrl *MockUsecase) StartCountingHashes(ctx context.Context, cfg *file.Config) error {
	if cfg == nil {
		return СonfigMissing
	}
	if ctrl.err != nil {
		return ctrl.err
	}
	return nil
}
