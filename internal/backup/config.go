package backup

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"rsync-backup/internal/filepaths"
	"rsync-backup/internal/slices"
	"rsync-backup/internal/types"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"

	resty "github.com/go-resty/resty/v2"
	"gopkg.in/ini.v1"
)

type Config struct {
	LocalIP       string `ini:"local_ip"`
	RsyncServer   string `ini:"rsync_server"`
	RsyncUser     string `ini:"rsync_user"`
	RsyncPassword string `ini:"rsync_password"`
	ConfigRootURL string `ini:"config_root_url"`
	Apps          []App
	DestModuleApp string
	DestModuleLog string
	RsyncCmd      string
	Today         string
}

type App struct {
	Name             string
	SourceDir        string `ini:"source_dir"`
	DestDir          string `ini:"dest_dir"`
	SourceFiles      string `ini:"source_files"`
	LogFiles         string `ini:"log_files"`
	SourceFileList   []string
	LogFileList      []string
	RsyncSourceDir   string
	RsyncDestDir     string
	RsyncSourceFiles []string
	RsyncLogFiles    []string
	DestAppURL       string
	DestLogURL       string
	RsyncPassword    string
}

var (
	Sections   []string
	ConfigFile string
	config     *Config
)

func GetRsyncCmd() string {
	sysType := runtime.GOOS
	rsyncCmd := ""
	rsyncPath := ""

	log.Printf("system type: %s", sysType)

	if rsyncPath == "" {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
			//return errors.Wrap(err, "os.Executable")
		}
		rsyncPath = filepath.Dir(ex)
	}

	if sysType == "windows" {
		rsyncCmd = filepath.Join(rsyncPath, "\\cwrsync\\bin\\rsync.exe")
	}

	if sysType == "linux" {
		rsyncCmd = "rsync"
	}

	log.Printf("rsync cmd: %s", rsyncCmd)

	return rsyncCmd
}

func (config *Config) GetLocalConfig(configFile string) *ini.File {

	ConfigFile = configFile
	cfg, err := ini.LoadSources(ini.LoadOptions{
		AllowPythonMultilineValues: true,
	}, configFile)

	if err != nil {
		log.Fatalf("读取配置文件%s失败: %v", configFile, err)
	}

	config.getConfig(cfg)

	return cfg
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

	Cfg, err := ini.Load(ini.LoadOptions{
		AllowPythonMultilineValues: true,
	}, configResult.Data.Content)

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

func (config *Config) getConfig(cfg *ini.File) {
	err := cfg.MapTo(config)

	if err != nil {
		log.Fatalf("配置文件内容不符合格式: %v", err)
	}

	config.Today = time.Now().Format("20060102")

	localIP := GetClientIP(config.ConfigRootURL + "/ip")

	if localIP != "" {
		config.LocalIP = localIP
	}

	config.DestModuleApp = "app_" + config.LocalIP
	config.DestModuleLog = "log_" + config.LocalIP
	config.RsyncCmd = GetRsyncCmd()

	Sections = cfg.SectionStrings()

	for _, section := range Sections {

		if section == "DEFAULT" {
			continue
		}

		log.Printf("section: %s", section)
		app := new(App)
		app.Name = section
		app.DestDir = "/"
		err = cfg.Section(section).MapTo(app)
		if err != nil {
			log.Fatalf("map to app error: %s", err)
		}

		r := regexp.MustCompile(`,|\s|\n|\r`)
		app.SourceFileList = slices.DeleteNil(r.Split(app.SourceFiles, -1))
		app.LogFileList = slices.DeleteNil(r.Split(app.LogFiles, -1))

		app.RsyncPassword = config.RsyncPassword
		app.DestAppURL = "rsync://" + config.RsyncUser + "@" +
			filepaths.DeleteExtraSeporator(config.RsyncServer+":/"+config.DestModuleApp+"/"+app.DestDir+"/"+config.Today+"/")

		app.DestLogURL = "rsync://" + config.RsyncUser + "@" +
			filepaths.DeleteExtraSeporator(config.RsyncServer+":/"+config.DestModuleLog+"/"+app.DestDir+"/")

		sourceDirBase := filepath.Base(app.SourceDir)

		if app.SourceDir != "" {
			app.RsyncSourceDir = filepaths.ConvertToRsyncPath(app.SourceDir, sourceDirBase)
			app.RsyncSourceFiles = append(app.RsyncSourceFiles, filepaths.GetDirFilesRsyncPath(app.SourceDir, sourceDirBase)...)
		}

		for _, sourceFile := range app.SourceFileList {
			app.RsyncSourceFiles = append(app.RsyncSourceFiles, filepaths.GetDirFilesRsyncPath(sourceFile, "/")...)
		}

		for _, logFilePattern := range app.LogFileList {
			if filepath.IsAbs(logFilePattern) {
				logFiles, _ := filepath.Glob(logFilePattern)

				for _, logFile := range logFiles {
					rsyncLogFile := filepaths.ConvertToRsyncPath(logFile, "/")
					app.RsyncLogFiles = append(app.RsyncLogFiles, rsyncLogFile)
				}

			} else {
				logFiles, _ := filepath.Glob(filepath.Join(app.SourceDir, logFilePattern))

				for _, logFile := range logFiles {
					rsyncLogFile := filepaths.ConvertToRsyncPath(logFile, sourceDirBase)
					app.RsyncLogFiles = append(app.RsyncLogFiles, rsyncLogFile)
				}
			}
		}

		app.RsyncSourceFiles = slices.DeleteSlice(app.RsyncSourceFiles, app.RsyncLogFiles)

		config.Apps = append(config.Apps, *app)
	}
}

func GetClientIP(URL string) string {
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
