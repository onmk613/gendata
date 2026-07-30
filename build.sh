#!/bin/bash

rm -rf bin
mkdir bin

gofmt -w .
go mod vendor

GOOS=linux  GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o bin/gendata-linux-amd64   .
GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o bin/gendata-darwin-arm64  .

rm -rf vendor
