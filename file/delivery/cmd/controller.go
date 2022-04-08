package cmd

import "github.com/avepa/gethash/file"

type Controller struct {
	ctrlUsecase file.Usecase
}

func NewController(ctrl file.Usecase) *Controller {
	return &Controller{
		ctrlUsecase: ctrl,
	}
}
