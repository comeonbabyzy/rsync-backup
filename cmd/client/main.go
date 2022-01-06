package main

import (
	"flag"
	"log"
	"rsync-backup/internal/backup"
)

var (
	h      bool
	config *backup.Config
)

func main() {

	configFile := ""
	configServer := ""

	flag.StringVar(&configFile, "c", "", "配置文件，默认为config.ini")
	flag.StringVar(&configServer, "s", "",
		"配置服务器地址，从此地址获取配置，默认配置URL为：http://ConfigServer:8080/config")

	flag.Parse()

	if h {
		flag.Usage()
	}

	if configFile == "" {
		configFile = "d:/Projects/rsync-backup/example/config.ini"
	}

	config = new(backup.Config)
	log.Printf("configFile: %s", configFile)

	if configFile != "" {
		config.GetLocalConfig(configFile)
	}

	for _, app := range config.Apps {
		app.BackupApp()
	}
}
