set GOOS=windows
set GOARCH=amd64
go build -o bin/win64/csdb-downloader-win64.exe
set GOOS=windows
set GOARCH=386
go build -o bin/win32/csdb-downloader-win32.exe
set GOOS=linux
set GOARCH=amd64
go build -o bin/lin64/csdb-downloader-amd64
set GOOS=linux
set GOARCH=arm64
go build -o bin/arm64/csdb-downloader-arm64
set GOOS=darwin
set GOARCH=amd64
go build -o bin/mac_amd64/csdb-downloader-mac-amd64
set GOOS=darwin
set GOARCH=arm64
go build -o bin/mac_arm64/csdb-downloader-mac-arm64