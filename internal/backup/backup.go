package backup

import (
	"fmt"
	"log"
	"rsync-backup/internal/commands"

	"rsync-backup/internal/filepaths"
)

func (app *App) BackupApp() error {

	rsyncFileName := "rsync_files.txt"
	rsyncLogName := "rsync_logs.txt"

	appFilesFrom := fmt.Sprintf(`--files-from=%s`, rsyncFileName)
	logFilesFrom := fmt.Sprintf(`--files-from=%s`, rsyncLogName)

	filepaths.WriteListToFile(rsyncFileName, app.RsyncSourceFiles)
	filepaths.WriteListToFile(rsyncLogName, app.RsyncLogFiles)

	rsyncCmd := "d:\\Projects\\rsync-backup\\bin\\cwrsync\\bin\\rsync.exe"

	resultApp, err := commands.RunCmd("RSYNC_PASSWORD="+app.RsyncPassword, rsyncCmd, "-avvvzRP", "--no-g",
		appFilesFrom, "/", app.DestAppURL)

	if err != nil {
		return err
	}

	log.Println(resultApp)

	resultLog, err := commands.RunCmd("RSYNC_PASSWORD="+app.RsyncPassword, rsyncCmd, "-avvvzRP", "--no-g",
		logFilesFrom, "/", app.DestLogURL)

	if err != nil {
		return err
	}

	log.Println(resultLog)

	return nil
}
