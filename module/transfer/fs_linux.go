package transfer

func IsDiskSpaceNotEnough(fileSize uint64) (bool, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(cfg.Get().Sys.Root, &stat)
	if err != nil {
		return false, err
	}

	//TODO check whether exactly match with arm fs
	reservedConf := cfg.Get().Sys.DiskReservePercent

	freeSpace := stat.Bfree * uint64(stat.Bsize)
	availSpace := stat.Bavail * uint64(stat.Bsize)

	if reservedConf == "-" {
		return availSpace <= fileSize, nil
	} else {
		reserved, err := strconv.ParseUint(reservedConf, 10, 64)
		if err != nil {
			return false, errors.Wrap(err)
		}

		reserveSpace := stat.Blocks * uint64(stat.Bsize) * reserved / 100
		return freeSpace <= fileSize+reserveSpace, nil
	}
}
