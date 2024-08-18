package transfer

import (
	"syscall"

	"github.com/madlabx/pkgx/errors"
	"github.com/madlabx/pkgx/log"

	"github.com/madlabx/fs/common/cfg"
)

func IsDiskSpaceNotEnough(fileSize uint64) (bool, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(cfg.Get().Sys.Root, &stat)
	if err != nil {
		return false, errors.Wrap(err)
	}

	freeSpace := stat.Bfree * uint64(stat.Bsize) / 1024
	availSpace := stat.Bavail * uint64(stat.Bsize) / 1024

	log.Errorf("freeSpace:%v", freeSpace)
	log.Errorf("availSpace:%v", availSpace)
	return freeSpace < fileSize, nil
}
