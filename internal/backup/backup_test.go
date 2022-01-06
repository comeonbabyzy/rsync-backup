package backup

import (
	"testing"
)

func TestBackup(t *testing.T) {
	configFile := "../../example/config.ini"

	config := new(Config)

	cfg := config.GetLocalConfig(configFile)

	t.Log(cfg)

	t.Log(config)

	config.Apps[0].BackupApp()
}
