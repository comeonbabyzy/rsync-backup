package filepaths

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func GetOwner(path string) (int, int, error) {
	return -1, -1, errors.New("Not implemented")
}

func Chown(rootPath string, uid, gid int) {
	log.Info("Not implemented")
}
