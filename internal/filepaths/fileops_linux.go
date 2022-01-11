package filepaths

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"syscall"
)

func GetOwner(path string) (int, int, error) {
	info, err := os.Stat(path)
	if err != nil {
		return -1, -1, err
	}
	stat := info.Sys().(*syscall.Stat_t)

	return int(stat.Uid), int(stat.Gid), nil
}

func Chown(rootPath string, uid, gid int) {
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		return os.Chown(path, uid, gid)
	})

	if err != nil {
		log.Infof("Chown %s fail. Error: %v", rootPath, err)
	}
}
