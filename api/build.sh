#!/bin/bash

cd "$(dirname "$0")"
set -e

if [[ " $* " == *" --help "* ]]; then
  echo "Использование: ./build.sh [--force] [--skip-js-build] [--skip-go-build]"
  echo "  --force         Перед началом сборки удалить папку build"
  echo "  --skip-js-build Пропустить сборку клиента nb-proxy JS API"
  echo "  --skip-go-build Пропустить сборку сервера nb-proxy Go API и клиента nb-back Go API"
  echo "Нужна установленная Java и wget"
  exit 0
fi

echo "Checking for required tools"
if ! command -v java &> /dev/null; then
  echo "Java is required!"
  exit 1
fi

if ! command -v wget &> /dev/null; then
  echo "wget is required!"
  exit 1
fi

if [[ " $* " == *" --force "* ]]; then
  echo "Forcing the build"
  rm -rf build
fi

echo "Preparing the build environment"
if [ ! -d "build" ]; then
  mkdir build
fi

echo "Downloading OpenAPI Generator CLI"
if [ ! -f "build/openapi-generator-cli.jar" ]; then
  wget https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/7.9.0/openapi-generator-cli-7.9.0.jar\
    -O build/openapi-generator-cli.jar
fi

if [[ " $* " == *" --skip-js-build "* ]]; then
  echo "Skipping the nb-proxy JS API client build"
else
  echo "Generating nb-proxy JS API client"
  java -jar ./build/openapi-generator-cli.jar generate\
    -i ./nb-proxy-api.yaml\
    -g javascript\
    -o ./build/nb-proxy-js\
    -c ./config/javascript.yaml
  (cd build/nb-proxy-js && npm install && npm run build)
fi

if [[ " $* " == *" --skip-go-build "* ]]; then
  echo "Skipping the nb-proxy Go API server build"
  echo "Skipping the nb-back Go API client build"
else
  echo "Generating nb-proxy Go API server"
  java -jar ./build/openapi-generator-cli.jar generate\
    -i ./nb-proxy-api.yaml\
    -g go-echo-server\
    -o ./build/nb-proxy-go\
    -c ./config/go-echo-server.yaml
  echo "Building nb-proxy Go API server"
  (cd build/nb-proxy-go && go mod tidy && go build)

  echo "Generating nb-back Go API client"
  java -jar ./build/openapi-generator-cli.jar generate\
    -i ./nb-back-api.yaml\
    -g go\
    -o ./build/nb-back-go\
    -c ./config/go.yaml
  echo "Building nb-back Go API client"
  (cd build/nb-back-go && go mod tidy && go build)
fi

echo "Build completed"