package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/avepa/gethash/file"
)

func TestSetError(t *testing.T) {
	mock := &MockUsecase{}
	err := errors.New("error")
	mock.SetError(err)
	if !errors.Is(err, mock.err) {
		t.Errorf("Invalid error. Expected: %v, received: %v.", err, mock.err)
	}
}

func TestMockStartCountingHashes(t *testing.T) {
	mock := &MockUsecase{}

	err := mock.StartCountingHashes(context.Background(), nil)
	if !errors.Is(err, СonfigMissing) {
		t.Errorf("Invalid error. Expected: %v, received: %v.", err, mock.err)
	}

	err = mock.StartCountingHashes(context.Background(), &file.Config{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	mock.SetError(errors.New("error"))
	err = mock.StartCountingHashes(context.Background(), &file.Config{})
	if errors.Is(err, errors.New("error")) {
		t.Errorf("Invalid error. Expected: %v, received: %v.", err, mock.err)
	}
}
