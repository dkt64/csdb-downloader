set GOOS=windows
set GOARCH=amd64
go build -o bin/csdb-downloader-win64/csdb-downloader-win64.exe
set GOOS=windows
set GOARCH=386
go build -o bin/csdb-downloader-win32/csdb-downloader-win32.exe
set GOOS=linux
set GOARCH=amd64
go build -o bin/csdb-downloader-linux_amd64/csdb-downloader-linux-amd64
set GOOS=linux
set GOARCH=arm64
go build -o bin/csdb-downloader-linux_arm64/csdb-downloader-linux-arm64
set GOOS=darwin
set GOARCH=amd64
go build -o bin/csdb-downloader-mac_amd64/csdb-downloader-mac-amd64
set GOOS=darwin
set GOARCH=arm64
go build -o bin/csdb-downloader-mac_arm64/csdb-downloader-mac-arm64