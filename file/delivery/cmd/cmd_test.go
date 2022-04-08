package cmd

import (
	"fmt"
	"os"
	"testing"
)

// func TestStart(t *testing.T) {
// 	ctrlUsecase := &usecase.MockUsecase{}
// 	ctrl := NewController(ctrlUsecase)
// 	err := ctrl.Start(context.Background())
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	// flag.CommandLine = flag.NewFlagSet("test", flag.ExitOnError)

// 	// flag.CommandLine.Usage = func() {
// 	// 	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
// 	// 	flag.PrintDefaults()
// 	// }
// 	// flag.Parse()

// 	// ctrlUsecase.SetError(errors.New("error"))
// 	// err = ctrl.Start(context.Background())
// 	// if err == nil {
// 	// 	t.Error(err)
// 	// }
// 	return
// }

func TestGetStartParams(t *testing.T) {
	flags := map[string]string{
		"max_proc":       "3",
		"folder_dir":     "main",
		"save_file_path": "main/hash",
		"save_file_name": "hash.txt",
	}

	for key, value := range flags {
		os.Args = append(os.Args, fmt.Sprintf("-%v=%v", key, value))
	}
	cfg, err := GetStartParams()
	if err != nil {
		t.Error(os.Args, cfg, err)
	}
	if cfg.MaxProc != 3 {
		t.Error("incorrect max_proc value")
	}
	if cfg.FolderDir != flags["folder_dir"] {
		t.Error("incorrect folder_dir value")
	}
	if cfg.SaveFileDir != flags["save_file_path"] {
		t.Error("incorrect save_file_path value")
	}
	if cfg.SaveFileName != flags["save_file_name"] {
		t.Error("incorrect save_file_name value")
	}
}
