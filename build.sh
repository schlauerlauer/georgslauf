#!/usr/bin/env bash

PACKAGE="georgslauf"
VERSION="$(git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null | sed 's/^.//')"
COMMIT_HASH="$(git rev-parse --short HEAD)"
BUILD_TIMESTAMP=$(date '+%Y-%m-%dT%H:%M:%S')

LDFLAGS=(
  "-X '${PACKAGE}/interfaces.version=${VERSION}'"
  "-X '${PACKAGE}/interfaces.commitHash=${COMMIT_HASH}'"
  "-X '${PACKAGE}/interfaces.buildTime=${BUILD_TIMESTAMP}'"
)

CGO_ENABLED=0 go build -ldflags="${LDFLAGS[*]}" -o ./tmp/main .
