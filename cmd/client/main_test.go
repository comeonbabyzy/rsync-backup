package main

import (
	"rsync-backup/internal/backup"
	"testing"
)

func TestMain(t *testing.T) {

	configFile := "d:/Projects/rsync-backup/example/config.ini"

	config := new(backup.Config)

	config.GetLocalConfig(configFile)

	config.Apps[0].BackupApp()
}
