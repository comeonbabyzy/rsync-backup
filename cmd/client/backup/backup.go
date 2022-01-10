package backup

import (
	"fmt"
	"rsync-backup/internal/commands"
	"runtime"

	"rsync-backup/internal/filepaths"
)

func (app *App) BackupApp() error {

	rsyncFileName := "rsync_files.txt"
	rsyncLogName := "rsync_logs.txt"

	appFilesFrom := fmt.Sprintf(`--files-from=%s`, rsyncFileName)
	logFilesFrom := fmt.Sprintf(`--files-from=%s`, rsyncLogName)
	rsyncCmd := "rsync"

	filepaths.WriteListToFile(rsyncFileName, app.RsyncSourceFiles)
	filepaths.WriteListToFile(rsyncLogName, app.RsyncLogFiles)

	if runtime.GOOS == "windows" {
		rsyncCmd = "c:\\temp\\cwrsync\\bin\\rsync.exe"
	}

	err := commands.RunCmd("RSYNC_PASSWORD="+app.RsyncPassword, rsyncCmd, "-avvzRP", "--no-g",
		appFilesFrom, "/", app.DestAppURL)

	if err != nil {
		return err
	}

	err = commands.RunCmd("RSYNC_PASSWORD="+app.RsyncPassword, rsyncCmd, "-avvzRP", "--no-g",
		logFilesFrom, "/", app.DestLogURL)

	if err != nil {
		return err
	}

	return nil
}
