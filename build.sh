#!/usr/bin/env bash

BIN_NAME=keyvent
VERSION=$(git describe --tags)

if [ $VERSION = "" ]; then
    VERSION="dev"
fi

printf "Compiling version ${VERSION} for linux/amd64 "

GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.Version=${VERSION}'" -o ./dist/${VERSION}/${BIN_NAME}
echo ✅

echo Output dir: ./dist/${VERSION}/
