package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/avepa/gethash/file"
	tHash "github.com/avepa/gethash/tools/hash"
)

type Config struct {
	MaxProc    int
	FolderPath string
	SaveFile   string
	HashName   string
}

func (ctrl *Controller) Start(ctx context.Context) error {
	cfg, err := GetStartParams()
	if err != nil {
		return nil
	}
	return ctrl.ctrlUsecase.StartCountingHashes(ctx, cfg)
}

func GetStartParams() (*file.Config, error) {
	hashName := flag.String(
		"hash",
		tHash.SHA3_256.String(),
		fmt.Sprintf("необходимо указать хеш-алгоритм. Доступны: %v.",
			strings.Join([]string{tHash.SHA3_256.String(),
				tHash.SHA3_512.String(),
				tHash.SHA256.String(),
				tHash.SHA512.String(),
				tHash.Keccak256.String(),
				tHash.Keccak512.String(),
				tHash.MD5.String(),
				tHash.Ripemd160.String(),
			},
				", ",
			),
		),
	)
	maxproc := flag.Int("max_proc", 2, "количество одновременно обрабатываемых файлов.")
	folderDir := flag.String("folder_dir", "", "Необходимо передать путь до папки.")
	saveFilePath := flag.String("save_file_path", "", "Необходимо передать путь куда сохранять файл. По умолчанию используется таже директория что и в folderPath")
	saveFileName := flag.String("save_file_name", "hashes_file.txt", "Необходимо указать имя создаваемого файла.")
	flag.Parse()

	if maxproc == nil || *maxproc < 1 {
		err := errors.New("maxproc cannot be lower than one")
		fmt.Println(err)
		return nil, err
	}
	if *saveFilePath == "" {
		saveFilePath = folderDir
	}

	cfg := &file.Config{
		MaxProc:      *maxproc,
		FolderDir:    *folderDir,
		SaveFileDir:  *saveFilePath,
		SaveFileName: *saveFileName,
		HashName:     *hashName,
	}

	return cfg, nil
}
