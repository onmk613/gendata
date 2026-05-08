#!/bin/bash

rm -rf bin
rm -rf vendor

mkdir bin

gofmt -w .
go mod vendor

GOOS=linux GOARCH=amd64 go build -o bin/gendata-linux-amd64 cmd/gendata/main.go
GOOS=darwin GOARCH=arm64 go build -o bin/gendata-darwin-arm64 cmd/gendata/main.go
