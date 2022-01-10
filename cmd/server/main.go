package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"rsync-backup/internal/filepaths"
	"rsync-backup/internal/types"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	rsyncDataDir                  string
	rsyncClientConfigFileTemplate string
	rsyncConfigFileTemplate       string
	rsyncConfigContentTemplate    string
}

var (
	serverConfig = ServerConfig{}
)

func getRemoteIP(c *gin.Context) string {
	remoteIP, _ := c.RemoteIP()
	return remoteIP.String()
}

func webRoot(c *gin.Context) {
	result := types.JsonResult{
		Code:    http.StatusOK,
		Message: "Hello, world",
		Data:    nil,
	}

	c.IndentedJSON(http.StatusOK, result)
}

func postServerConfig(c *gin.Context) {
	var result types.ResponseBase
	var returnStatus int
	var returnMessage string

	ip := getRemoteIP(c)

	rootPath := "/volume1/databackup"
	appPath := rootPath + "/app" + "/" + ip
	logPath := rootPath + "/log" + "/" + ip

	rsyncConfigFile := fmt.Sprintf(serverConfig.rsyncConfigFileTemplate, ip)

	os.MkdirAll(appPath, os.ModePerm)
	os.MkdirAll(logPath, os.ModePerm)

	content := fmt.Sprintf(serverConfig.rsyncConfigContentTemplate, ip, appPath, ip, ip, logPath, ip)
	err := filepaths.WriteStringToFile(rsyncConfigFile, content)

	if err != nil {
		returnStatus = http.StatusBadRequest
		returnMessage = fmt.Sprintf("写入配置文件失败 %v", err)

	} else {
		returnStatus = http.StatusOK
		returnMessage = fmt.Sprintf("写入配置文件成功 %s", rsyncConfigFile)
	}

	result = types.ResponseBase{
		Code:    returnStatus,
		Message: returnMessage,
	}
	c.IndentedJSON(returnStatus, result)
	log.Println(result)
}

func getConfig(c *gin.Context) {

	var result types.ResponseGetConfig

	rsyncClientConfigFile := fmt.Sprintf(serverConfig.rsyncClientConfigFileTemplate, getRemoteIP(c))

	content, err := filepaths.ReadTxtFile(rsyncClientConfigFile)

	if err != nil {
		log.WithFields(log.Fields{
			"file":      rsyncClientConfigFile,
			"remote_ip": getRemoteIP(c),
		}).Error("读取配置文件失败")

		result = types.ResponseGetConfig{
			ResponseBase: types.ResponseBase{
				Code:    http.StatusBadRequest,
				Message: "Fail",
			},
			Data: types.DataGetConfig{
				FileName: "",
				Content:  "",
			},
		}
		c.IndentedJSON(http.StatusBadRequest, result)
		return
	}

	log.WithFields(log.Fields{
		"file":      rsyncClientConfigFile,
		"remote_ip": getRemoteIP(c),
	}).Info("读取配置文件成功")

	result = types.ResponseGetConfig{
		ResponseBase: types.ResponseBase{
			Code:    http.StatusOK,
			Message: "Success",
		},
		Data: types.DataGetConfig{
			FileName: rsyncClientConfigFile,
			Content:  content,
		},
	}

	c.IndentedJSON(http.StatusOK, result)
}

func getClientIP(c *gin.Context) {
	ResultIP := types.ResponseGetIP{
		ResponseBase: types.ResponseBase{
			Code:    http.StatusOK,
			Message: "Success",
		},
		Data: types.DataGetIP{
			IP: getRemoteIP(c),
		},
	}
	c.IndentedJSON(http.StatusOK, ResultIP)
}

func postConfig(c *gin.Context) {
	requestCfg := types.RequestPostConfig{}
	c.BindJSON(&requestCfg)
	rsyncClientConfigFile := fmt.Sprintf(serverConfig.rsyncClientConfigFileTemplate, getRemoteIP(c))
	err := filepaths.WriteStringToFile(rsyncClientConfigFile, requestCfg.Data.Content)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest,
			types.ResponseBase{
				Code:    http.StatusBadRequest,
				Message: "写入文件失败！",
			})
	}

	c.IndentedJSON(http.StatusOK,
		types.ResponseBase{
			Code:    http.StatusOK,
			Message: "写入文件成功！",
		})
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/", webRoot)
	r.GET("/ip", getClientIP)
	r.GET("/config", getConfig)
	r.POST("/config", postConfig)
	r.POST("/serverconfig", postServerConfig)
	r.GET("/cwrsync.zip", func(c *gin.Context) {
		c.File("cwrsync.zip")
	})
	return r
}

func (serverConfig *ServerConfig) getConfig() {

	flag.StringVar(&serverConfig.rsyncDataDir, "d", "", "备份数据根目录")
	flag.Parse()

	if serverConfig.rsyncDataDir == "" {
		flag.Usage()
		log.Fatal("请加 -d 参数设置备份数据根目录")
	}

	var err error
	serverConfig.rsyncDataDir, err = filepath.Abs(serverConfig.rsyncDataDir)

	if err != nil {
		log.Fatal(err)
	}

	serverConfig.rsyncClientConfigFileTemplate = serverConfig.rsyncDataDir + "/%s.conf"
	serverConfig.rsyncConfigFileTemplate = "/etc/rsyncd.d/%s.conf"
	serverConfig.rsyncConfigContentTemplate = `
[app_%s]
path = %s
uid = rsync
gid = users 
read only = no
auth users = rsync
secrets file = /etc/rsyncd.secrets
hosts allow = %s

[log_%s]
path = %s
uid = rsync
gid = users 
read only = no
auth users = rsync
secrets file = /etc/rsyncd.secrets
hosts allow = %s
`
}

func main() {

	serverConfig.getConfig()

	r := setupRouter()
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
