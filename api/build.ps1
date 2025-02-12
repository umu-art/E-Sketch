Set-Location (Split-Path -Path $MyInvocation.MyCommand.Definition -Parent)
$ErrorActionPreference = "Stop"

Function Show-Help {
    Write-Host "Использование: .\build.ps1 [--force] [--skip-js-build] [--skip-go-build]"
    Write-Host "  --force         Перед началом сборки удалить папку build"
    Write-Host "  --skip-js-build Пропустить сборку клиента est-proxy JS API"
    Write-Host "  --skip-go-build Пропустить сборку сервера est-proxy Go API и клиента est-back Go API"
    Write-Host "Нужна установленная Java и wget"
    Exit 0
}

if ($args -contains "--help") {
    Show-Help
}

Write-Host "Checking for required tools"
if (-not (Get-Command java -ErrorAction SilentlyContinue)) {
    Write-Host "Java is required!"
    Exit 1
}

if (-not (Get-Command wget -ErrorAction SilentlyContinue)) {
    Write-Host "wget is required!"
    Exit 1
}

if ($args -contains "--force") {
    Write-Host "Forcing the build"
    Remove-Item -Recurse -Force -ErrorAction SilentlyContinue -Path .\build
}

Write-Host "Preparing the build environment"
if (-not (Test-Path -Path .\build)) {
    New-Item -ItemType Directory -Path .\build | Out-Null
}

Write-Host "Downloading OpenAPI Generator CLI"
if (-not (Test-Path -Path .\build\openapi-generator-cli.jar)) {
    Invoke-WebRequest -Uri "https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/7.9.0/openapi-generator-cli-7.9.0.jar" -OutFile .\build\openapi-generator-cli.jar
}

if ($args -contains "--skip-js-build") {
    Write-Host "Skipping the est-proxy JS API client build"
} else {
    Write-Host "Generating est-proxy JS API client"
    java -jar ./build/openapi-generator-cli.jar generate -i ./est-proxy-api.yaml -g javascript -o ./build\est-proxy-js -c ./config/javascript.yaml
    Set-Location -Path .\build\est-proxy-js
    npm install
    npm run build
    Set-Location -Path (Split-Path -Path $MyInvocation.MyCommand.Definition -Parent)
}

if ($args -contains "--skip-go-build") {
    Write-Host "Skipping the est-proxy Go API server build"
    Write-Host "Skipping the est-back Go API client build"
} else {
    Write-Host "Generating est-proxy Go API server"
    Start-Process java -ArgumentList "-jar ./build/openapi-generator-cli.jar generate -i ./est-proxy-api.yaml -g go-echo-server -o ./build/est-proxy-go -c ./config/go-echo-server.yaml" -Wait
    Write-Host "Building est-proxy Go API server"
    Set-Location -Path .\build\est-proxy-go
    go mod tidy
    go build
    Set-Location -Path (Split-Path -Path $MyInvocation.MyCommand.Definition -Parent)
    
    Write-Host "Generating est-back Go API client"
    Start-Process java -ArgumentList "-jar ./build/openapi-generator-cli.jar generate -i ./est-back-api.yaml -g go -o ./build/est-back-go -c ./config/go.yaml" -Wait
    Write-Host "Building est-back Go API client"
    Set-Location -Path .\build\est-back-go
    go mod tidy
    go build
    Set-Location -Path (Split-Path -Path $MyInvocation.MyCommand.Definition -Parent)
}

Write-Host "Build completed"
