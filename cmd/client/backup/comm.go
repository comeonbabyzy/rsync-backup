package backup

import (
	"io/ioutil"
	"net/http"
	"rsync-backup/internal/types"

	resty "github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

func GetClientIP(URL string) string {
	log.Printf("get client ip URL: %s", URL)
	client := resty.New()
	result := types.ResponseGetIP{}

	resp, err := client.R().SetResult(&result).Get(URL)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	if resp.StatusCode() != http.StatusOK {
		log.Fatal(resp.Status())
		return ""
	}

	return result.Data.IP
}

func GetClientIPPlain(URL string) string {

	resp, err := http.Get(URL)
	if err != nil {
		log.Printf("get client ip error: %v", err)
		return ""
	}
	defer resp.Body.Close()
	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("get client ip error: %v", err)
		return ""
	}
	return string(s)
}

func MakeServerConfig(URL string) {
	client := resty.New()

	resp, err := client.R().Post(URL)

	if err != nil {
		log.Error(err)
		return
	}

	if resp.StatusCode() != http.StatusOK {
		log.Error(err)
		return
	}

	log.Info("save config success")

}

func GetCWRsync(URL string, destFile string) {
	client := resty.New()
	_, err := client.R().
		SetOutput(destFile).
		Get(URL)

	if err != nil {
		log.Error(err, "get cwrsync error")
		return
	}
}

func (config *Config) GetServerConfig(URL string) {

	client := resty.New()
	configResult := &types.ResponseGetConfig{}
	resp, err := client.R().SetResult(configResult).Get(URL)

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode() != 200 {
		log.Fatal(resp.Status())
	}

	Cfg, err := ini.LoadSources(ini.LoadOptions{
		AllowPythonMultilineValues: true,
	}, []byte(configResult.Data.Content))

	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}
	config.getConfig(Cfg)
}

func (config *Config) SaveServerConfig(URL string, content string) {
	client := resty.New()
	body := types.RequestPostConfig{
		Data: types.DataPostConfig{
			Content: content,
		},
	}
	resp, err := client.R().SetBody(body).Post(URL)

	if err != nil {
		log.Error(err)
		return
	}

	if resp.StatusCode() != http.StatusOK {
		log.Error(err)
		return
	}

	log.Info("save config success")
}
