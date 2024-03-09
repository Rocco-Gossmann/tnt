#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o tnt.linux.x64
GOOS=linux GOARCH=arm64 go build -o tnt.linux.arm64
GOOS=windows GOARCH=amd64 go build -o tnt.windows.x64.exe
GOOS=windows GOARCH=386 go build -o tnt.windows.x32.exe
GOOS=darwin GOARCH=amd64 go build -o tnt.mac.x64
GOOS=darwin GOARCH=arm64 go build -o tnt.mac.arm64
