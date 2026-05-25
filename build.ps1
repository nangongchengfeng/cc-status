# CC Status 构建脚本 (Windows PowerShell)

param(
    [ValidateSet("all", "web", "server", "clean")]
    [string]$Target = "all"
)

$ErrorActionPreference = "Stop"
$ScriptDir = $PSScriptRoot

function Build-Web {
    Write-Host "Building web..." -ForegroundColor Cyan
    Push-Location "$ScriptDir\web"
    npm run build
    Pop-Location

    Write-Host "Copying web dist to server..." -ForegroundColor Cyan
    Remove-Item -Recurse -Force "$ScriptDir\server\internal\handler\ui\dist" -ErrorAction SilentlyContinue
    New-Item -ItemType Directory -Force "$ScriptDir\server\internal\handler\ui\dist" | Out-Null
    Copy-Item -Recurse -Force "$ScriptDir\web\dist\*" "$ScriptDir\server\internal\handler\ui\dist\"
}

function Build-Server {
    Write-Host "Building server (statically linked)..." -ForegroundColor Cyan
    Push-Location "$ScriptDir\server"
    $env:CGO_ENABLED = 0
    $exeName = if ($env:OS -eq "Windows_NT") { "bin\server.exe" } else { "bin/server" }
    go build -tags embed -ldflags="-s -w" -o $exeName ./cmd/server
    Pop-Location
}

function Clean {
    Write-Host "Cleaning..." -ForegroundColor Cyan
    Remove-Item -Recurse -Force "$ScriptDir\server\bin" -ErrorAction SilentlyContinue
    Remove-Item -Recurse -Force "$ScriptDir\server\internal\handler\ui\dist" -ErrorAction SilentlyContinue
    New-Item -ItemType Directory -Force "$ScriptDir\server\internal\handler\ui\dist" | Out-Null
    New-Item -ItemType File -Force "$ScriptDir\server\internal\handler\ui\dist\.gitkeep" | Out-Null
    Remove-Item -Recurse -Force "$ScriptDir\web\dist" -ErrorAction SilentlyContinue
}

switch ($Target) {
    "web" { Build-Web }
    "server" { Build-Server }
    "clean" { Clean }
    "all" { Clean; Build-Web; Build-Server }
}

Write-Host "Build completed!" -ForegroundColor Green
