#!/usr/bin/env bash

PACKAGE="georgslauf"
VERSION="$(git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null | sed 's/^.//')"
BUILD_TIMESTAMP=$(date '+%Y-%m-%dT%H:%M:%S')

LDFLAGS=(
  "-X '${PACKAGE}/handler.version=${VERSION}'"
  "-X '${PACKAGE}/handler.buildTime=${BUILD_TIMESTAMP}'"
)

CGO_ENABLED=1 go build -ldflags="${LDFLAGS[*]}" -o ./tmp/georgslauf .
