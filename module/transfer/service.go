package transfer

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/madlabx/pkgx/errors"
	"github.com/madlabx/pkgx/log"
	"go.uber.org/zap"

	"github.com/madlabx/fs/common/cfg"
)

func doClean(root string, f func(string) error) error {
	fileList, err := os.ReadDir(root)
	if err != nil {
		return errors.Wrap(err)
	}

	for _, dir := range fileList {
		if !dir.IsDir() {
			log.Errorf("Found invalid file, do clean manually, file:%v", dir.Name())
			continue
		}

		if dir.Name() == "0" {
			//skip 0
			continue
		}

		expireData := 0
		expireData, err = strconv.Atoi(dir.Name())
		if err != nil {
			log.Errorf("Found wrong dir, name:%v", dir.Name())
			continue
		}

		err = filepath.Walk(dir.Name(), func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				log.Errorf("walk dir:%v failed, path:%v, err:%v", dir.Name(), path, err)
				return nil
			}

			if fi == nil {
				return err
			}
			if fi.IsDir() {
				if path == dir.Name() {
					return nil
				}

				fileInDir, err := os.ReadDir(path)
				if err != nil {
					log.Errorf("read dir failed, path:%v, err:%v", path, err)
					return nil
				}
				if len(fileInDir) == 0 {
					zap.L().Info("empty dir", zap.String("path", path))
					if err = f(path); err != nil {
						log.Errorf("failed to do f(path), path:%v, err:%v", path, err)
					}
				}
				return nil
			}
			fileTime := fi.ModTime()
			if time.Since(fileTime) > time.Hour*24*time.Duration(expireData) {
				log.Infof("file:%v has expired", path)
				if err = f(path); err != nil {
					log.Errorf("failed to do f(path), path:%v, err:%v", path, err)
				}
			}

			return nil
		})
	}

	return err
}

func Launch(ctx context.Context) error {
	go func() {
		root := cfg.Get().Sys.Root

		if root == "" || root == "/" {
			return
		}

		ticker := time.NewTicker(time.Second * 60)
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				if err := doClean(root, func(path string) (err error) {
					err = os.RemoveAll(path)
					if err != nil {
						log.Errorf("remove dir:%v, has error:%v", path, err)
					}
					return
				}); err != nil {
					log.Errorf("Igore doClean err:%v", err)
				}
			}
		}
	}()

	return nil
}
