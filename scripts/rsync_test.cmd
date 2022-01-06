set RSYNC_PASSWORD=test
 ..\build\cwrsync\bin\rsync.exe -avvvzRP --no-g --files-from=rsync_files.txt --exclude-from=rsync_logs.txt / rsync://rsync@192.168.191.143:/app_192.168.191.133/

rem ..\build\cwrsync\bin\rsync.exe -avvvzRP --no-g --files-from=rsync-files.txt --exclude-from=log-files.txt / /cygdrive/d/temp/

rem ..\build\cwrsync\bin\rsync.exe -avvvzRP --no-g --exclude-from=log-files.txt /cygdrive/d/./test/backuptest rsync://rsync@192.168.191.143:/app_192.168.191.133/

rem ..\build\cwrsync\bin\rsync.exe -avvvzRP --no-g --exclude-from=log-files.txt /cygdrive/d/./rsync_app_test rsync://rsync@192.168.191.143:/app_192.168.191.133/