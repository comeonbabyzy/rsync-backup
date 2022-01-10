package filepaths

import (
	"archive/zip"
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

func GetAbsURL(baseURL string, uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		log.Fatal(err)
	}
	base, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	return base.ResolveReference(u).String()
}

func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func ReadTxtFile(path string) (string, error) {
	txt, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.Wrap(err, "读取文件失败。")
	}

	content := string(txt)
	return content, nil
}

func WriteStringToFile(filename string, data string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(data))
	if err != nil {
		return err
	}

	return nil
}

func WriteListToFile(fileName string, strList []string) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

	if err != nil {
		//return errors.Wrap(err, "open file "+fileName+" error.")
		panic(err)
	}

	datawriter := bufio.NewWriter(file)

	for _, str := range strList {
		_, _ = datawriter.WriteString(str + "\n")
	}

	datawriter.Flush()
	defer file.Close()
}

// windows 下将 c:\windows 改成 /cygdrive/c/windows格式，用于rsync命令进行同步
func ToCygwinPath_OLD(path string) string {

	if runtime.GOOS == "windows" {
		volumeName := strings.ToLower(filepath.VolumeName(path))
		drive := strings.ReplaceAll(volumeName, ":", "")

		return strings.ReplaceAll(filepath.ToSlash(path), volumeName, "/cygdrive/"+drive)
	}

	return path
}

// windows 下将 c:\windows 改成 /cygdrive/c/windows格式，用于rsync命令进行同步
// 方法：删除冒号（:），前加/cygdrive，并将\替换成/
// 使用前应将路径转换为绝对路径
func ToCygwinPath(path string) string {
	if runtime.GOOS == "windows" {
		noColonPath := strings.ReplaceAll(path, ":", "")
		return strings.ToLower("/cygdrive/" + filepath.ToSlash(noColonPath))
	}

	return path
}

//将路径加/./，用于rsync命令进行同步，以指定目的目录
//如： 模式1：指定最后一级目录 /usr/local/nginx -> /usr/local/./nginx，目的目录则为/nginx
//或： 模式2：指定全目录：c:\windows\system32\ -> /./c/windows/system32, rsync目的目录则为/c/windows/system32
// 使用时会将路径设置为绝对路径
func ToRsyncSourcePath(sourcePath string, destPath string) string {
	absSourcePath, _ := filepath.Abs(sourcePath)
	destSlashPath := filepath.ToSlash(destPath)
	rsyncSourcePath := ""
	if destSlashPath == "/" {
		rsyncSourcePath = "/./" + filepath.ToSlash(absSourcePath)
	} else {
		rsyncSourcePath = strings.Replace(filepath.ToSlash(sourcePath), filepath.ToSlash(destSlashPath), "/./"+destSlashPath, 1)
	}
	return rsyncSourcePath
}

func DeleteExtraSeporator(path string) string {
	splashPath := strings.ReplaceAll(path, "\\", "/")
	reg := regexp.MustCompile(`/+`)

	returnPath := reg.ReplaceAllString(splashPath, "/")

	return returnPath

}

func ConvertToRsyncPath(path string, destPath string) string {
	return DeleteExtraSeporator(ToCygwinPath(ToRsyncSourcePath(path, destPath)))
}

func GetDirFilesRsyncPath(path string, destPath string) []string {
	var files []string

	if IsFile(path) {
		return []string{ConvertToRsyncPath(path, destPath)}
	}

	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		absPath, _ := filepath.Abs(path)
		files = append(files, ConvertToRsyncPath(absPath, destPath))
		//log.Printf("%s\n", f.Name())
		return nil
	})

	if err != nil {
		log.Printf("filepath.Walk() returned %v\n", err)
	}

	return files
}

//参考 https://www.jianshu.com/p/4593cfffb9e9
func Unzip(zipFile string, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		fpath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
