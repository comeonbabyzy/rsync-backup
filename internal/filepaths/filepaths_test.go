package filepaths

import (
	"fmt"
	"path/filepath"

	"testing"
)

func TestGetDirFiles(t *testing.T) {
	path := "..\\"
	files := GetDirFilesRsyncPath(path, "/")

	t.Log(files)
}

func TestFilePath(t *testing.T) {
	path := "d:\\test"
	// path = "logs\\*.*"
	// path = "d:\\test\\*.*"
	path2 := "d:\\test\\backuptest\\logs\\*.*"

	absPath, _ := filepath.Abs(path)
	base := filepath.Base(absPath)
	dir := filepath.Dir(absPath)
	toSlash := filepath.ToSlash(absPath)
	allExeFiles, _ := filepath.Glob(path + "*.exe")
	isAbs := filepath.IsAbs(path)
	relPath, _ := filepath.Rel(path, path2)

	fmt.Printf("Abs path %s \n", absPath)
	fmt.Printf("Base path %s \n", base)
	fmt.Printf("Dir path %s \n", dir)
	fmt.Printf("All exe files %s \n", allExeFiles)
	fmt.Printf("To slash path %s \n", toSlash)
	fmt.Printf("Is abs path %t \n", isAbs)
	fmt.Printf("Rel path %s \n", relPath)
	fmt.Println(path)

	t.Log(path)
}

func TestConvertToRsyncPath1(t *testing.T) {
	rsyncPath1, rsyncPath2, rsyncPath3 := ConvertToRsyncPath2("/usr/local/nginx")

	t.Logf("Linux: rsyncPath1: %s, rsyncPath2: %s, rsyncPath3: %s \n", rsyncPath1, rsyncPath2, rsyncPath3)

	winRsyncPath1, winRsyncPath2, winRsyncPath3 := ConvertToRsyncPath2("d:\\test\\")

	t.Logf("Windows: winRsyncPath1: %s, winRsyncPath2: %s, winRsyncPath3: %s \n", winRsyncPath1, winRsyncPath2, winRsyncPath3)
}

func TestConvertToRsyncPath2(t *testing.T) {

	path := "d:\\test"

	//cygwinPath := util.ToCygwinPath(path)
	rsyncPath := ToRsyncSourcePath(path, `/test`)
	cygwinPath := ToCygwinPath(rsyncPath)

	t.Logf("Rsync path: %s, %s \n", rsyncPath, cygwinPath)

}
