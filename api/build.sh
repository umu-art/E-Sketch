#!/bin/bash

cd "$(dirname "$0")"
set -e

if [[ " $* " == *" --help "* ]]; then
  echo "Использование: ./build.sh [--force] [--skip-js-build] [--skip-go-build]"
  echo "  --force         Перед началом сборки удалить папку build"
  echo "  --skip-js-build Пропустить сборку клиента est-proxy JS API"
  echo "  --skip-go-build Пропустить сборку сервера est-proxy Go API и клиента est-back Go API"
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
declare -a pids

if [[ " $* " == *" --skip-js-build "* ]]; then
  echo "Skipping the est-proxy JS API client build"
else
  {
    echo "Generating est-proxy JS API client"
    java -jar ./build/openapi-generator-cli.jar generate\
      -i ./est-proxy-api.yaml\
      -g javascript\
      -o ./build/est-proxy-js\
      -c ./config/javascript.yaml
    (cd build/est-proxy-js && npm install && npm run build)
    echo "est-proxy JS API client build completed"
  } &
  pids+=($!)
fi

if [[ " $* " == *" --skip-go-build "* ]]; then
  echo "Skipping the est-proxy Go API server build"
  echo "Skipping the est-back Go API client build"
else
  {
    echo "Generating est-proxy Go API server"
    java -jar ./build/openapi-generator-cli.jar generate\
      -i ./est-proxy-api.yaml\
      -g go-echo-server\
      -o ./build/est-proxy-go\
      -c ./config/go-echo-server.yaml
    echo "Building est-proxy Go API server"
    (cd build/est-proxy-go && go mod tidy && go build)
    echo "Generating est-back Go API client"
    java -jar ./build/openapi-generator-cli.jar generate\
      -i ./est-back-api.yaml\
      -g go\
      -o ./build/est-back-go\
      -c ./config/go.yaml
    echo "Building est-back Go API client"
    (cd build/est-back-go && go mod tidy && go build)
    echo "est-proxy and est-back Go API builds completed"
  } &
  pids+=($!)
fi

if [[ " $* " == *" --skip-cpp-build "* ]]; then
  echo "Skipping the est-back-cpp C++ API server build"
else
  {
    echo "Generating est-back-cpp C++ API server"
    java -jar ./build/openapi-generator-cli.jar generate\
      -i ./est-back-api.yaml\
      -g cpp-pistache-server\
      -o ./build/est-back-cpp\
      -c ./config/cpp.yaml
    echo "est-back-cpp C++ API server build completed"
  } &
  pids+=($!)
fi

if [[ " $* " == *" --skip-java-build "* ]]; then
  echo "Skipping the est-mono-api API server build"
else
  {
    echo "Generating est-mono-api API server"
    java -jar ./build/openapi-generator-cli.jar generate\
      -i ./est-proxy-api.yaml\
      -g spring\
      -o ./build/est-mono-api\
      -c ./config/java-spring.yaml
    echo "Building est-mono-api API server"
    (cd build/est-mono-api && mvn clean install)
    echo "est-mono-api API server build completed"
  } &
  pids+=($!)
fi

for pid in "${pids[@]}"; do
  wait "$pid"
done

echo "Build completed"