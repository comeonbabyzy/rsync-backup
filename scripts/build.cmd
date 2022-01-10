@echo off

SET CGO_ENABLED=0  
SET GOARCH=amd64
SET GOOS=windows

go build -ldflags="-H windows" -o bin\rsync_backup.exe cmd\client\main.go
go build -o bin\rsync_backup_server.exe cmd\server\main.go

SET CGO_ENABLED=0  
SET GOOS=linux
SET GOARCH=amd64

go build -o bin\rsync_backup cmd\client\main.go
go build -o bin\rsync_backup_server cmd\server\main.go

copy example\config.ini bin\config.ini
