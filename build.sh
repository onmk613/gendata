#!/bin/bash

rm -rf bin
rm -rf vendor

mkdir bin

gofmt -w .
go mod vendor

GOOS=linux GOARCH=amd64 go build -o bin/gendata-linux-amd64 ./main.go
GOOS=darwin GOARCH=arm64 go build -o bin/gendata-darwin-arm64 ./main.go
