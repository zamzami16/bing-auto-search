$ErrorActionPreference = "Stop"
$BIN = "bin"
New-Item -Force -ItemType Directory $BIN | Out-Null

$env:CGO_ENABLED = "0"
$env:GOOS = "windows"
$env:GOARCH = "amd64"

# Auto-detect version from git tag or use 'dev'
$VERSION = "dev"
try {
    $gitTag = git describe --tags --exact-match 2>$null
    if ($LASTEXITCODE -eq 0 -and $gitTag) {
        $VERSION = $gitTag
    }
} catch {
    # Ignore error, use dev
}

Write-Host "Building version: $VERSION"
$LDFLAGS = "-s -w -X main.version=$VERSION"

go build -trimpath -ldflags "$LDFLAGS" -o "$BIN\desktop-bing-auto.exe" .\cmd\desktop-bing-auto
Write-Host "Built -> $BIN\desktop-bing-auto.exe ($VERSION)"
