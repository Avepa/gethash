package file

import "context"

type Usecase interface {
	StartCountingHashes(ctx context.Context, cfg *Config) error
}

type Config struct {
	FolderDir    string
	SaveFileDir  string
	SaveFileName string
	HashName     string
	MaxProc      int
}
