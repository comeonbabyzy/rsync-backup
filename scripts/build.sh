
export CGO_ENABLED=0  
export GOARCH=amd64
export GOOS=windows

go build -ldflags="-H windows" -o bin/rsync_backup.exe cmd/client/main.go
go build -o bin/rsync_backup_server.exe cmd/server/main.go

export CGO_ENABLED=0  
export GOOS=linux
export GOARCH=amd64

go build -o bin/rsync_backup cmd/client/main.go
go build -o bin/rsync_backup_server cmd/server/main.go

cp example/config.ini bin/config.ini
