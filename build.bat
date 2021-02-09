set GOOS=windows
set GOARCH=amd64
go build -o bin/win64/csdb-downloader.exe
set GOOS=windows
set GOARCH=386
go build -o bin/win32/csdb-downloader.exe
set GOOS=linux
set GOARCH=amd64
go build -o bin/lin64/csdb-downloader
set GOOS=linux
set GOARCH=arm64
go build -o bin/arm64/csdb-downloader
set GOOS=darwin
set GOARCH=amd64
go build -o bin/mac64/csdb-downloader