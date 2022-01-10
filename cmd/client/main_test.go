package main

import (
	"rsync-backup/cmd/client/backup"
	"testing"
)

func TestMain1(t *testing.T) {

	configFile := "d:/Projects/rsync-backup/example/config.ini"

	config := new(backup.Config)

	config.GetLocalConfig(configFile)

	config.Apps[0].BackupApp()
}

func TestMain(t *testing.T) {
	main()
}
