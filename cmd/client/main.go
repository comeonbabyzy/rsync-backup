package main

import (
	"flag"
	"rsync-backup/cmd/client/backup"
	"rsync-backup/internal/filepaths"

	log "github.com/sirupsen/logrus"
)

var (
	h      bool
	config *backup.Config
)

func initLog(level log.Level) {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	log.SetLevel(level)
}

func main() {

	initLog(log.InfoLevel)
	configFile := ""
	configServerIP := ""
	configServerPort := ""

	flag.StringVar(&configFile, "c", "", "配置文件，默认为config.ini")
	flag.StringVar(&configServerIP, "s", "",
		"配置服务器IP，如192.168.191.143，从此地址获取配置，默认配置ROOT URL则为：http://192.168.191.143:8080")
	flag.StringVar(&configServerPort, "p", "8080",
		"配置服务器端口，如8080，从此端口获取配置，默认配置ROOT URL则为：http://192.168.191.143:8080")

	flag.Parse()

	if h {
		flag.Usage()
	}
	/*
		if configFile == "" {
			configFile = "d:/Projects/rsync-backup/example/config.ini"
		}
	*/
	config = new(backup.Config)

	log.WithFields(log.Fields{
		"configfile": configFile,
	}).Info("读取本地配置文件。")

	if configFile != "" {
		config.GetLocalConfig(configFile)

		if config.ConfigRootURL != "" {
			content, err := filepaths.ReadTxtFile(configFile)
			if err != nil {
				log.Error("读取配置文件失败。")
			}

			config.SaveServerConfig(config.ConfigRootURL+"/config", content)
		}
	} else {
		if configServerIP != "" {
			configRootURL := "http://" + configServerIP + ":" + configServerPort
			config.GetServerConfig(configRootURL + "/config")
		}
	}

	backup.GetCWRsync(
		config.ConfigRootURL+"/cwrsync.zip", "c:\\temp\\cwrsync.zip")

	filepaths.Unzip("c:\\temp\\cwrsync.zip", "c:\\temp")

	backup.MakeServerConfig(
		config.ConfigRootURL + "/serverconfig")

	for _, app := range config.Apps {
		app.BackupApp()
	}
}
