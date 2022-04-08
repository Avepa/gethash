package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/avepa/gethash/file/delivery/cmd"
	"github.com/avepa/gethash/file/repository/os"
	"github.com/avepa/gethash/file/usecase"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	ctrlRepo := os.NewController()
	ctrlUsecase := usecase.NewController(ctrlRepo)
	ctrlDeliv := cmd.NewController(ctrlUsecase)

	err := ctrlDeliv.Start(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
}
