#!/bin/bash

set -euo pipefail
cd "$(dirname "$0")"

rm -rf bin
mkdir bin

gofmt -w main.go cmd internal
go mod vendor
trap 'rm -rf vendor' EXIT

GOOS=linux  GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o bin/gendata-linux-amd64   .
GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o bin/gendata-darwin-arm64  .
