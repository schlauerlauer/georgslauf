#!/usr/bin/env bash

PACKAGE="georgslauf"
VERSION="$(git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null | sed 's/^.//')"
BUILD_TIMESTAMP=$(date '+%Y-%m-%dT%H:%M:%S')

LDFLAGS=(
	"-X '${PACKAGE}/internal/handler.version=${VERSION}'"
	"-X '${PACKAGE}/internal/handler.buildTime=${BUILD_TIMESTAMP}'"
	"-s"
	"-w"
	'-extldflags "-static"'
)

# suggested on https://github.com/mattn/go-sqlite3/tree/master?tab=readme-ov-file#arm
CGO_ENABLED=1 CC=arm-linux-gnueabihf-gcc CXX=arm-linux-gnueabihf-g++ GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="${LDFLAGS[*]}" -o /app/build/georgslauf -v .

# CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64 go build -ldflags="${LDFLAGS[*]}" -o /app/build/georgslauf -v .
