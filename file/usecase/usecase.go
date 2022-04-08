package usecase

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"golang.org/x/exp/slices"

	"github.com/avepa/gethash/file"
)

func (ctrl *Controller) StartCountingHashes(ctx context.Context, cfg *file.Config) error {
	cntrlFH, err := newControllerFileHash(cfg.MaxProc, cfg.HashName, ctrl.ctrlRepo)
	if err != nil {
		fmt.Println("failed to get config, err:", err)
		return err
	}

	files, err := ctrl.ctrlRepo.GetAllNameFileList(cfg.FolderDir)
	if err != nil {
		fmt.Println("failed to get list of files, err:", err)
		return err
	}
	slices.Sort(files)

	cntrlFH.wgSaveHash.Add(1)
	var sfCtx context.Context
	sfCtx, cntrlFH.shCancel = context.WithCancel(ctx)
	err = cntrlFH.saveHash(sfCtx, cfg.SaveFileDir, cfg.SaveFileName, files)
	if err != nil {
		fmt.Println("failed to create file, err:", err)
		return err
	}

	var chCtx context.Context
	chCtx, cntrlFH.chCancel = context.WithCancel(ctx)
	cntrlFH.wgCountHash.Add(cfg.MaxProc)
	for i := 0; i < cfg.MaxProc; i++ {
		go cntrlFH.countHash(chCtx, cfg.FolderDir)
	}

	for _, name := range files {
		cntrlFH.nameFile <- name
	}

	ticker := time.NewTicker(time.Second / 2)
	for _ = range ticker.C {
		if len(cntrlFH.nameFile) == 0 {
			ticker.Stop()
			cntrlFH.StopCountHash()
			break
		}
	}
	cntrlFH.StopSaveHash()
	return nil
}

func (ctrl *ControllerFileHash) StopCountHash() {
	ctrl.chCancel()
	ctrl.wgCountHash.Wait()
	close(ctrl.nameFile)
	return
}

func (ctrl *ControllerFileHash) StopSaveHash() {
	ctrl.shCancel()
	ctrl.wgSaveHash.Wait()
	close(ctrl.hashFile)
	return
}

func (ctrl *ControllerFileHash) countHash(ctx context.Context, folderDir string) {
	for {
		select {
		case <-ctx.Done():
			ctrl.wgCountHash.Done()
			return
		case nameFile := <-ctrl.nameFile:
			f, err := ctrl.ctrlRepo.GetFile(folderDir, nameFile)
			if err != nil {
				ctrl.hashFile <- HashFile{
					FileName: nameFile,
					Error:    err,
				}
				continue
			}

			hasher := ctrl.newHash()
			if _, err = io.Copy(hasher, f); err != nil {
				ctrl.hashFile <- HashFile{
					FileName: nameFile,
					Error:    err,
				}
				f.Close()
				continue
			}

			ctrl.hashFile <- HashFile{
				FileName: nameFile,
				FileHash: hasher.Sum(nil),
			}
			f.Close()
		}
	}
}

func (ctrl *ControllerFileHash) saveHash(ctx context.Context, saveFileDir, saveFileName string, filesName []string) (err error) {
	file, err := ctrl.ctrlRepo.CreateFile(saveFileDir, saveFileName)
	if err != nil {
		return
	}

	go func() {
		countedFile := make(map[string]string)
		defer file.Close()
		for {
			select {
			case <-ctx.Done():
				ctrl.wgSaveHash.Done()
				return
			case hf := <-ctrl.hashFile:
				if len(filesName) > 0 {
					if filesName[0] == hf.FileName {
						var data string
						if hf.Error != nil {
							data = fmt.Sprintf("Error: %v\n", hf.Error)
						} else {
							data = fmt.Sprintf("%v\n", hex.EncodeToString(hf.FileHash))
						}
						_, err := file.WriteString(data)
						if err != nil {
							fmt.Println(err)
						}

						filesName = filesName[1:]
						for _, fn := range filesName {
							if data, ok := countedFile[fn]; ok {
								_, err := file.WriteString(data)
								if err != nil {
									fmt.Println("Failed to save data to file, err:", err)
								}

								filesName = filesName[1:]
								delete(countedFile, fn)
							} else {
								break
							}
						}
					} else {
						var data string
						if hf.Error != nil {
							data = fmt.Sprintf("Error: %v\n", hf.Error)
						} else {
							data = fmt.Sprintf("%v\n", hex.EncodeToString(hf.FileHash))
						}
						countedFile[hf.FileName] = data
						continue
					}
				}
				if len(filesName) == 0 {
					ctrl.wgSaveHash.Done()
					return
				}
			}
		}
	}()
	return
}
