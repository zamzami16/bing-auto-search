$ErrorActionPreference = "Stop"
$BIN = "bin"
New-Item -Force -ItemType Directory $BIN | Out-Null

$env:CGO_ENABLED = "0"
$env:GOOS = "windows"
$env:GOARCH = "amd64"

go build -trimpath -ldflags "-s -w" -o "$BIN\desktop-bing-auto.exe" .\cmd\desktop-bing-auto
Write-Host "Built -> $BIN\desktop-bing-auto.exe"
